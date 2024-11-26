package config

type Options struct {
	Namespace       string
	Threadiness     int
	HTTPSListenPort int

	CloudweavControllerUsername string
	GarbageCollectionUsername   string
}
