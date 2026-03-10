// PostRepositoryインターフェース。
// 「何ができるか」だけを定義し、「どうやるか」はpersistence層に任せます。
// これにより、テスト時にモック(偽物)に差し替えることも容易になります。
package domain

// PostRepository: 投稿の永続化操作を抽象化したインターフェースです。
type PostRepository interface {
	GetAll() ([]Post, error)
	GetByID(id int) (Post, error)
	Create(title, content, author string) (Post, error)
	Delete(id int) (bool, error)
}
