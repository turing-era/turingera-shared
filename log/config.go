package log

import (
	"os"
	"time"

	rotateLogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Levels zapcore level
var Levels = map[string]zapcore.Level{
	"":      zapcore.DebugLevel,
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
}

// InitLog 初始化日志
func InitLog() {
	var cores []zapcore.Core
	writer := viper.GetString("log.writer")
	switch writer {
	case "console":
		cores = append(cores, newConsoleCore())
	case "file":
		cores = append(cores, newFileCore())
	default:
		panic("log writer invalid")
	}

	logger := zap.New(
		zapcore.NewTee(cores...),
		zap.AddCallerSkip(1),
		zap.AddCaller(),
	)
	sugar = logger.Sugar()
}

// DefaultTimeFormat 默认时间格式
func DefaultTimeFormat(t time.Time) []byte {
	t = t.Local()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	micros := t.Nanosecond() / 1000

	buf := make([]byte, 23)
	buf[0] = byte((year/1000)%10) + '0'
	buf[1] = byte((year/100)%10) + '0'
	buf[2] = byte((year/10)%10) + '0'
	buf[3] = byte(year%10) + '0'
	buf[4] = '-'
	buf[5] = byte((month)/10) + '0'
	buf[6] = byte((month)%10) + '0'
	buf[7] = '-'
	buf[8] = byte((day)/10) + '0'
	buf[9] = byte((day)%10) + '0'
	buf[10] = ' '
	buf[11] = byte((hour)/10) + '0'
	buf[12] = byte((hour)%10) + '0'
	buf[13] = ':'
	buf[14] = byte((minute)/10) + '0'
	buf[15] = byte((minute)%10) + '0'
	buf[16] = ':'
	buf[17] = byte((second)/10) + '0'
	buf[18] = byte((second)%10) + '0'
	buf[19] = '.'
	buf[20] = byte((micros/100000)%10) + '0'
	buf[21] = byte((micros/10000)%10) + '0'
	buf[22] = byte((micros/1000)%10) + '0'
	return buf
}

func newEncoder() zapcore.Encoder {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "N",
		CallerKey:     "C",
		MessageKey:    "M",
		StacktraceKey: "S",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendByteString(DefaultTimeFormat(t))
		},
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	return encoder
}

func newConsoleCore() zapcore.Core {
	lvl := zap.NewAtomicLevelAt(Levels[viper.GetString("log.level")])
	return zapcore.NewCore(
		newEncoder(),
		zapcore.Lock(os.Stdout),
		lvl)
}

func newFileCore() zapcore.Core {
	var ws zapcore.WriteSyncer
	WithMaxAge := viper.GetDuration("log.with_max_age")
	WithRotationTime := viper.GetDuration("log.with_rotation_time")
	WithRotationCount := viper.GetInt("log.with_rotation_count")
	LogName := viper.GetString("log.name")

	writer, err := rotateLogs.New(
		LogName+".%Y%m%d%H%M",
		rotateLogs.WithLinkName(LogName),                // 生成软链，指向最新日志文件
		rotateLogs.WithMaxAge(WithMaxAge),               // 文件最大保存时间
		rotateLogs.WithRotationTime(WithRotationTime),   // 日志切割时间间隔
		rotateLogs.WithRotationCount(WithRotationCount), // 设置文件清理前最多保存的个数。
	)
	if err != nil {
		panic(err)
	}
	ws = zapcore.AddSync(writer)
	// 日志级别
	lvl := zap.NewAtomicLevelAt(Levels[viper.GetString("log.level")])
	return zapcore.NewCore(
		newEncoder(),
		ws, lvl,
	)
}
