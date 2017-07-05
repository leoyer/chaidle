package handler

type HandlerProcessor interface {
	DoProcessor(interface{})
	AddHanler(handler Handler)
}

func NewDefaultHanlerProcessor() HandlerProcessor{
	return  &defaultHandlerProcessor{}
}

type defaultHandlerProcessor struct {
	handlers []Handler
}

func (processor *defaultHandlerProcessor)AddHanler(handler Handler){
	processor.handlers = append(processor.handlers,handler)
}


func (processor *defaultHandlerProcessor)DoProcessor(obj interface{}){
	for _,handler := range processor.handlers{
		if handler.Handler(obj) != nil{
			return
		}
	}
}




