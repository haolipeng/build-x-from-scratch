package main

// RegexMatcher 实现基于正则表达式的匹配器
type RegexMatcher struct {
	pattern *Pattern
}

// NewRegexMatcher 创建一个新的正则表达式匹配器
func NewRegexMatcher(pattern *Pattern) *RegexMatcher {
	return &RegexMatcher{pattern: pattern}
}

// Match 判断文本中是否包含匹配模式的子串
// 这是 grep 的默认行为：在文本中查找匹配的子串
func (m *RegexMatcher) Match(text string) bool {
	// 尝试从文本的每个位置开始匹配
	for i := 0; i <= len(text); i++ {
		if m.matchAt(text, i) {
			return true
		}
	}
	return false
}

// MatchFull 判断整个文本是否完全匹配模式
// 用于需要精确匹配的场景
func (m *RegexMatcher) MatchFull(text string) bool {
	return m.matchAt(text, 0) && m.consumesAll(text)
}

// matchAt 从文本的指定位置开始尝试匹配模式
func (m *RegexMatcher) matchAt(text string, startPos int) bool {
	return m.matchTokens(m.pattern.tokens, text, 0, startPos)
}

// matchTokens 递归匹配 token 序列
// patIdx: 当前要匹配的 pattern token 索引
// textIdx: 当前要匹配的 text 字符索引
func (m *RegexMatcher) matchTokens(tokens []Token, text string, patIdx int, textIdx int) bool {
	// 基础情况 1: 模式已全部匹配
	if patIdx >= len(tokens) {
		return true
	}

	// 基础情况 2: 文本耗尽但模式未完
	if textIdx >= len(text) {
		return false
	}

	// 获取当前要匹配的 token
	token := tokens[patIdx]

	// 根据 token 类型进行匹配
	switch token.Type {
	case TokenWildcard:
		// 通配符 . 匹配任意单个字符
		return m.matchTokens(tokens, text, patIdx+1, textIdx+1)

	case TokenLiteral:
		// 字面字符必须精确匹配
		if text[textIdx] == token.Value {
			return m.matchTokens(tokens, text, patIdx+1, textIdx+1)
		}
		return false

	default:
		// 未知的 token 类型
		return false
	}
}

// consumesAll 检查从位置 0 开始匹配是否消耗了整个文本
func (m *RegexMatcher) consumesAll(text string) bool {
	// 这个方法在完整匹配时使用
	// 简化实现：检查模式长度是否等于文本长度
	return len(m.pattern.tokens) == len(text)
}

// Pattern 返回匹配器使用的模式
func (m *RegexMatcher) Pattern() string {
	return m.pattern.Raw()
}

// FindMatch 在文本中查找第一个匹配的子串
// 返回匹配的起始位置和结束位置（-1 表示未找到）
func (m *RegexMatcher) FindMatch(text string) (start, end int) {
	for i := 0; i <= len(text); i++ {
		if m.matchAt(text, i) {
			// 找到匹配，计算结束位置
			end := i + len(m.pattern.tokens)
			if end <= len(text) {
				return i, end
			}
		}
	}
	return -1, -1
}
