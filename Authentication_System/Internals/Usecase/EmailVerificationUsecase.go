package usecase

import (
	"context"
	"log"
	"time"
)

type EmailVerificationUsecase struct {
	verRepo EmailVerificationTokenRepoitory
	tokenGenerator TokenProvider
}

func NewEmailVerificationUsecase(vr EmailVerificationTokenRepoitory , tokenGenerator TokenProvider) *EmailVerificationUsecase  {
	  return &EmailVerificationUsecase{verRepo: vr , tokenGenerator: tokenGenerator}
}


type VerficationInput struct {
	Token string
}

type VerificationOutput struct {
	AccessToken  string    `json:"accessToken"`  
	RefreshToken string    `json:"refreshToken"` 
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (ev *EmailVerificationUsecase) Execute(ctx context.Context , input VerficationInput) (*VerificationOutput ,error) {
	token , err := ev.verRepo.FindByToken(ctx , input.Token);
	if err != nil {
		return  &VerificationOutput{}, err
	}

	if err := token.ValidateExpiration(); err != nil {
		return &VerificationOutput{} ,err
	}
    
     
	user , err := ev.verRepo.FindByID(ctx , token.UserID);
	if err != nil {
		return &VerificationOutput{}, err
	}

	 if err := user.MarkAsVerified(); err != nil {
	     if delErr := ev.verRepo.Delete(ctx , token.ID); delErr != nil{
			log.Printf("[WARN] Failed to delete token %s after verification failure: %v", token.ID, delErr)
		 }
		 return nil, err
	 }



	 if err := ev.verRepo.Save(ctx , user); err != nil {
		return &VerificationOutput{},err
	 }
     
	 accessToken , refreshToken , err := ev.tokenGenerator.GenerateTokenPair(user.ID() , user.Role());
	 if err != nil {
		return &VerificationOutput{},err
	 }
      

     return &VerificationOutput{
          AccessToken: accessToken,
		  RefreshToken: refreshToken,
		  Role: string(user.Role()),
		  CreatedAt: time.Now().UTC(),
		  UpdatedAt: time.Now().UTC(),
	 }, nil
}  