package validator

type MovieValidator struct {
	Title       string `validate:"required"`
	Description string `validate:"required"`
	Duration    string `validate:"required"`
	Artist      string `validate:"required"`
	Genre       string `validate:"required"`
	VideoUrl    string `validate:"required"`
}
