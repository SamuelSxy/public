package biz

type Tran interface {
}

type User struct {
	uid  int32
	name string
}

type UserRepo interface {
	SaveUser(*User)
	Begin() Tran
}

type UserUsecase struct {
	repo UserRepo
}

func (uc *UserUsecase) NewUserUsecase(u *User) error {

	tr := uc.repo.Begin()
	err := uc.repo.SaveUser(u)

	if err != nil {
		tr.Roollback()
	}

	return tr.Commit()

}
