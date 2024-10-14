package models

type Notification int

const (
	ErrorNotification Notification = iota
	SuccessNotification
	InfoNotification
)

func (n Notification) String() string {
	switch n {
	case SuccessNotification:
		return "success"
	case InfoNotification:
		return "info"
	default:
		return "error"
	}
}

type Toast struct {
	Message string
	Level   Notification
}

func NewToastNotification(message string, level Notification) Toast {
	return Toast{
		Message: message,
		Level:   level,
	}
}

// func (t Toast) Html() string {
// 	return fmt.Sprintf(`
//         <div hx-on::load="setTimeout(() => { this.remove() }, 2000)" id="toaster"
//         class="fixed z-50 bottom-5 right-5 rounded-md w-56 p-5 border %s">
//         <p class="font-medium">%s</p>
//         </div>
//         `, t.LevelColor(), t.Message)
// }
//
// func (t Toast) LevelColor() string {
// 	switch t.Level {
// 	case ErrorNotification:
// 		return "border-rose-400"
// 	case SuccessNotification:
// 		return "border-emerald-500"
// 	case InfoNotification:
// 		return "border-sky-400"
// 	}
//
// 	return "border-gray-300"
// }
