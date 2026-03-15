package clients

type BankClient struct{

}

func NewBankClient() BankClient{
	return BankClient{

	}
}

func (b* BankClient) Debit(accountId string, amount int64) error {
	return nil
}

func (b* BankClient) Credit(accountId string, amount int64) error {
	return nil
}
