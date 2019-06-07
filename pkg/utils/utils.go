package utils

//type EndPointStruct struct {
//	Name     string
//	Function eprequest.EpFunction
//}
//
//var EndPoints = []EndPointStruct{
//	{Name: "noise-spectrogram", Function: mustang.Get},
//}
//
//func getEndPointFunction(epname string) (eprequest.EpFunction, error) {
//
//	var epfunc eprequest.EpFunction = nil
//	found := false
//	for _, ep := range EndPoints{
//		if ep.Name == epname {
//			epfunc = ep.Function
//			found = true
//			break
//		}
//	}
//	if !found {
//		return epfunc, fmt.Errorf("unsupported endpoint name %s", epname)
//	}
//
//	return epfunc, nil
//}

//func IrisWSGet(req *eprequest.EPRequest) (*epresult.EPResult, error) {
//
//	resbuf, err := mustang.Get(req)
//	return &epresult.EPResult{req, resbuf}, err
//}
