package request

//func apiQueryPost(uri string, params ApiParams, secrets binanceStructs.Secrets) (interface{}, error) {
//	return apiQuery(uri, params, secrets, "POST",
//		func(params ApiParams, apiEndpoint string, uri string) (string, *bytes.Buffer, error) {
//			postContent, err := json.Marshal(params)
//			if err != nil {
//				return "", nil, err
//			}
//			buffer := bytes.NewBuffer(postContent)
//
//			var builder strings.Builder
//			builder.WriteString(apiEndpoint)
//			builder.WriteString(uri)
//			builder.WriteByte('?')
//			requestUrl := builder.String()
//
//			return requestUrl, buffer, nil
//		})
//}
