package repository

func checkBalance(balance, amount int) error {
	if balance < amount {
		return ErrNotEnoughCoins
	}
	return nil
}
