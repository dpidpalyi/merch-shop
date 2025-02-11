package service

type UserService struct {
	userRepo *repository.UserRepository
}

func (u *UserService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, userRepo.ErrRecordNotFound) {
			jwt, err := u.Add(username, password)
			return jwt, err
		}
		return "", nil
	}

	if !utils.CheckPasswordHash(user.password, password) {
		return "", Err
	}
}

func (u *UserService) Add(ctx context.Context, username, password string) (string, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	user := &models.User{
		Username: username,
		Password: hashedPassword,
	}

	token, err := utils.GetJWT(user)

	err = u.userRepo.Add(ctx, user)
	if err != nil {
		return "", err
	}
}
