package service

import (
	"zrsf.com/hbase_client/types"
	"zrsf.com/hbase_client/helper"
	"github.com/asaskevich/govalidator"
	"errors"
	"zrsf.com/hbase_client/handler"
	"zrsf.com/hbase_client/handler/save"
)

type HSaveService interface{
	Save(request types.SaveRequest) (error,string)
	SavePdf(request types.SaveRequest)  (error,string)
	SaveInvoiceIndexs(request types.SaveRequest)  (error,string)
}


type hsaveService struct {

}


func NewHsaveService() HSaveService{
	return &hsaveService{}
}


func (hsaveService) Save(request types.SaveRequest) (error,string ){
	if flag , _ := govalidator.ValidateStruct(request); !flag {
		return errors.New("XML内参数错误"),helper.CODE_CSERROR
	}
	return  helper.NewTsunaGoHbase().SaveWithRowKey(request),""
}
func (hsaveService)SavePdf(request types.SaveRequest)  (error,string){
	if flag , _ := govalidator.ValidateStruct(request); !flag {
		return errors.New("XML内参数错误"),helper.CODE_CSERROR
	}
	return  helper.NewTsunaGoHbase().SaveWithRowKey(request),""
}
func (hsaveService)SaveInvoiceIndexs(request types.SaveRequest)  (error,string){
	if flag , _ := govalidator.ValidateStruct(request); !flag {
		return errors.New("XML内参数错误"),helper.CODE_CSERROR
	}

	processor := handler.NewDefaultHanlerProcessor()

	processor.AddHanler(&save.PackIndexsParamHandler{})
	processor.AddHanler(&save.PackFPTAccountHandler{})
	processor.AddHanler(&save.InsertIntoMajorTableHandler{})
	processor.AddHanler(&save.OpenInvoiceCompanyIndexHandler{})
	processor.AddHanler(&save.ReceiveInvoiceIndexHandler{})
	processor.AddHanler(&save.FPTAccountIndexHandler{})
	processor.AddHanler(&save.MobilePhoneIndexHandler{})
	processor.AddHanler(&save.EmailIndexHandler{})



	processor.AddHanler(&save.PackRedEinvoiceParamHandler{})
	processor.AddHanler(&save.RedEinvoiceHandler{})

	processor.DoProcessor(request)


	return  nil,""
}