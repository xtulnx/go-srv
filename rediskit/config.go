package rediskit

type RedisConfig struct {
	Network  string // tcp
	Address  string //'127.0.0.1:6379'
	Password string // =''
	Db       int    // ='7'
}
