package models

type UserInterface interface{
	FindUserInfo(userId string)(Session,error)
}

type DbUserModel struct{}

func (userModel *DbUserModel)FindUserInfo(userId string)(Session,error){
	return Session{},nil
}