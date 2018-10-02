package request

import (
	"book-collections-restapi/dto/response"
	"book-collections-restapi/helpers"
	"github.com/go-playground/locales/en_US"
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en2 "gopkg.in/go-playground/validator.v9/translations/en"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type BookRequest struct {
	ID bson.ObjectId `json:"id"`
	Title string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
	Publisher string `json:"publisher" validate:"required"`
}

var (
	validate *validator.Validate
	uni *ut.UniversalTranslator
)

func (bookRequest *BookRequest) Validate(w http.ResponseWriter, r *http.Request) bool {
	en := en_US.New()
	uni = ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	validate = validator.New()
	en2.RegisterDefaultTranslations(validate, trans)

	var idValidation error

	if r.Method == "PUT" {
		idValidation = validate.Var(bookRequest.ID, "required")
	}

	err := validate.Struct(bookRequest)
	if err != nil || idValidation != nil {
		errs := response.ValidationResponse{}
		errs.Error = "Bad Request."
		errs.Message = "You have an error validation."

		var errorItem []string

		if idValidation != nil {
			errorItem = append(errorItem, "Id" + idValidation.(validator.ValidationErrors)[0].Translate(trans))
		}

		if err != nil {
			for _, e := range err.(validator.ValidationErrors) {
				errorItem = append(errorItem, e.Translate(trans))
			}
		}

		errs.Errors = errorItem

		helpers.RespondWithValidationError(w, errs)

		return true
	}

	return false
}