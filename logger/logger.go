package logger

import (
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/ini.v1"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 声明一个全局的SugaredLogger变量
var (
	// SugarLogger 全局的 SugarLogger 和 Log 实例
	SugarLogger *zap.SugaredLogger
	Log         *log.Logger
)

// ZapLogWriter 是一个适配器,用于将标准log的输出重定向到zap
type ZapLogWriter struct {
	logger *zap.Logger
}

// 进行时间的格式化定义
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// NewCustomTimeEncoderConfig 进行zap输出格式的自定义
func NewCustomTimeEncoderConfig() zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "time"       // 将时间戳字段名称设置为 "time"
	config.MessageKey = "message" //将msg改为"event"
	config.CallerKey = "code"     //将caller改为"code"
	config.EncodeTime = customTimeEncoder
	return config
}

// Write 实现io.Writer接口
func (z *ZapLogWriter) Write(p []byte) (n int, err error) {
	// 使用zap记录日志
	z.logger.Info(string(p))
	return len(p), nil
}

func Init() {
	// 加载配置文件
	cfg, err := ini.Load("server.conf")
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}
	// 获取DebugMode
	debugMode := cfg.Section("Server").Key("DebugMode").String()
	logFilePath := "./edgetts.log"

	// 使用lumberjack作为zap的日志写入器,实现日志轮滚
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFilePath, // 日志文件路径
		MaxSize:    10,          // 每个日志文件保存的最大尺寸 单位:MB
		MaxBackups: 10,          // 日志文件最多保存多少个备份
		MaxAge:     10,          // 文件最多保存多少天
		Compress:   false,       // 是否压缩
	}
	writeSyncer := zapcore.AddSync(lumberjackLogger)

	// 配置zap的日志格式
	var atomicLevel zap.AtomicLevel
	switch debugMode {
	case "debug":
		atomicLevel = zap.NewAtomicLevelAt(zap.DebugLevel) //调试信息输出
	case "info":
		atomicLevel = zap.NewAtomicLevelAt(zap.InfoLevel) //一般信息输出
	case "warn":
		atomicLevel = zap.NewAtomicLevelAt(zap.WarnLevel) //警告级别
	case "error":
		atomicLevel = zap.NewAtomicLevelAt(zap.ErrorLevel) //错误级别
	default:
		atomicLevel = zap.NewAtomicLevelAt(zap.InfoLevel) // 默认级别为Info
	}
	// 修改这里来使用自定义的时间格式
	// 使用自定义的时间编码配置
	encoderConfig := NewCustomTimeEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 初始化zap日志器
	core := zapcore.NewCore(encoder, writeSyncer, atomicLevel)
	logger := zap.New(core, zap.AddCaller(), zap.Fields(zap.String("name", "TTSProxy")))

	// 创建SugaredLogger
	SugarLogger = logger.Sugar()

	// 配置标准log的前缀和标志
	log.SetPrefix("[edgetts] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 配置log模块使用同一个日志文件
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("打开文件失败: %v", err)
	}
	Log = log.New(logFile, "[edgetts] ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
}

// ZapLoggerAdapter 是一个结构体,它包含了一个指向 zap.SugaredLogger 的指针.
// 这个结构体被用来将标准的日志消息适配到 zap logger.
type ZapLoggerAdapter struct {
	sugarLogger *zap.SugaredLogger
}

// NewZapLoggerAdapter 是 ZapLoggerAdapter 的构造函数.它接受一个指向 zap.SugaredLogger 的指针
// 并返回一个指向新创建的 ZapLoggerAdapter 的指针.
func NewZapLoggerAdapter(sugarLogger *zap.SugaredLogger) *ZapLoggerAdapter {
	return &ZapLoggerAdapter{sugarLogger: sugarLogger}
}

// Write 方法使 ZapLoggerAdapter 满足 io.Writer 接口.这个方法将传入的字节切片转换为字符串,
// 然后使用 zap logger 进行记录.这是必要的,因为许多 Go 的标准库和第三方库期望一个 io.Writer
// 来进行日志记录,而这个方法允许 ZapLoggerAdapter 充当这个角色.
func (z *ZapLoggerAdapter) Write(p []byte) (n int, err error) {
	// Here we convert the byte slice to a string and log it
	z.sugarLogger.Debugln(string(p))
	return len(p), nil
}
