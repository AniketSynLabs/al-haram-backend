package handler

func ok() map[string]string               { return map[string]string{"status": "ok"} }
func err400(err error) map[string]string  { return map[string]string{"error": err.Error()} }
func err500(err error) map[string]string  { return map[string]string{"error": err.Error()} }
