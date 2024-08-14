package Infrastructure

func CompareHashAndPasswordCustom(hashedPassword string, plainPassword string) bool {
	// err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	// return err == nil
	if plainPassword == hashedPassword {
		return true
	}
	return false
}

func GenerateFromPasswordCustom(password string) (string, error) {
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// return string(hashedPassword), err
	return password, nil
}

// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
