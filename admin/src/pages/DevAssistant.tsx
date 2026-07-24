import { useState, useEffect, useRef } from 'react';
import './DevAssistant.scss';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeHighlight from 'rehype-highlight';
import 'highlight.js/styles/github-dark.css';
import {
  Layout,
  Button,
  Input,
  Card,
  Modal,
  List,
  Tag,
  Spin,
  Empty,
  Tooltip,
} from 'antd';
import {
  MessageOutlined,
  PlusOutlined,
  DeleteOutlined,
  SendOutlined,
  RobotOutlined,
  UserOutlined,
  FileTextOutlined,
  CodeOutlined,
  BulbOutlined,
  BugOutlined,
  DatabaseOutlined,
  BarChartOutlined,
  EyeOutlined,
  TeamOutlined,
} from '@ant-design/icons';
import type { ChatSession, ChatMessage } from '@/types';
import {
  getSessions,
  createSession,
  getSessionDetail,
  deleteSession,
  confirmAction,
  chatWithAssistant,
} from '@/api';

const { Sider, Content } = Layout;
const { TextArea } = Input;

const quickCommands = [
  { label: '修复这个 Bug', icon: <BugOutlined style={{ fontSize: 14 }} />, prompt: '请描述你遇到的 Bug，我来帮你分析和修复。' },
  { label: '生成 CRUD 接口', icon: <DatabaseOutlined style={{ fontSize: 14 }} />, prompt: '请告诉我需要创建的实体名称和字段，我来生成完整的 CRUD 接口。' },
  { label: '优化这段代码', icon: <CodeOutlined style={{ fontSize: 14 }} />, prompt: '请提供需要优化的代码，我来分析并给出优化建议。' },
  { label: '添加新功能', icon: <PlusOutlined style={{ fontSize: 14 }} />, prompt: '请描述你想要添加的新功能，我来帮你实现。' },
  { label: '生成图表', icon: <BarChartOutlined style={{ fontSize: 14 }} />, prompt: '请提供数据和图表类型，我来生成可视化图表。' },
  { label: '代码审查', icon: <EyeOutlined style={{ fontSize: 14 }} />, prompt: '请提供需要审查的代码，我来进行代码审查。' },
  { label: '接口调试', icon: <TeamOutlined style={{ fontSize: 14 }} />, prompt: '请提供接口地址和参数，我来帮你调试。' },
  { label: '读取文件', icon: <CodeOutlined style={{ fontSize: 14 }} />, prompt: '请提供文件路径，我来帮你读取文件内容。' },
];

export const DevAssistant = () => {
  const [sessions, setSessions] = useState<ChatSession[]>([]);
  const [currentSession, setCurrentSession] = useState<ChatSession | null>(null);
  const [inputValue, setInputValue] = useState('');
  const [loading, setLoading] = useState(false);
  const [confirmModalVisible, setConfirmModalVisible] = useState(false);
  const [createSessionModalVisible, setCreateSessionModalVisible] = useState(false);
  const [newSessionTitle, setNewSessionTitle] = useState('');
  const [scrollBottom, setScrollBottom] = useState<HTMLDivElement | null>(null);

  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    loadSessions();
  }, []);

  useEffect(() => {
    if (scrollBottom) {
      scrollBottom.scrollIntoView({ behavior: 'smooth' });
    }
  }, [currentSession, scrollBottom]);

  const loadSessions = async () => {
    try {
      const res = await getSessions();
      if (res.data.sessions) {
        setSessions(res.data.sessions);
        if (res.data.sessions.length > 0 && !currentSession) {
          setCurrentSession(res.data.sessions[0]);
        }
      }
    } catch (error) {
      console.error('加载会话失败', error);
    }
  };

  const handleCreateSession = async () => {
    try {
      setNewSessionTitle('');
      setCreateSessionModalVisible(true);
    } catch (error) {
      console.error('创建会话失败', error);
    }
  };

  const handleConfirmCreateSession = async () => {
    try {
      const res = await createSession(newSessionTitle.trim());
      if (res.data.session) {
        setSessions((prev) => [res.data.session, ...prev]);
        setCurrentSession(res.data.session);
        setInputValue('');
      }
      setCreateSessionModalVisible(false);
    } catch (error) {
      console.error('创建会话失败', error);
    }
  };

  const handleSelectSession = async (sessionId: string) => {
    const session = sessions.find((s) => s.id === sessionId);
    if (session) {
      try {
        const res = await getSessionDetail(sessionId);
        setCurrentSession(res.data.session);
      } catch (error) {
        console.error('获取会话详情失败', error);
        setCurrentSession(session);
      }
    }
  };

  const handleDeleteSession = (sessionId: string) => {
    Modal.confirm({
      title: '删除会话',
      content: '确定要删除这个会话吗？删除后无法恢复。',
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await deleteSession(sessionId);
          setSessions((prev) => prev.filter((s) => s.id !== sessionId));
          if (currentSession?.id === sessionId) {
            setCurrentSession(null);
          }
        } catch (error) {
          console.error('删除会话失败', error);
        }
      },
    });
  };

  const handleSendMessage = async () => {
    if (!inputValue.trim() || !currentSession) return;

    const message = inputValue.trim();
    setInputValue('');
    setLoading(true);

    const newMessages: ChatMessage[] = [
      ...(currentSession?.messages || []),
      {
        id: Date.now().toString(),
        role: 'user',
        content: message,
        created_at: new Date().toISOString(),
      },
    ];

    setCurrentSession((prev) => prev ? { ...prev, messages: newMessages } : null);

    try {
      const response = await chatWithAssistant({
        session_id: currentSession.id,
        message,
      });
      if (response.data.session) {
        setCurrentSession(response.data.session);
      }

      setLoading(false);
    } catch (error: any) {
      console.error('发送消息失败', error);
      Modal.error({
        title: '发送失败',
        content: error.msg || error.message || '发送消息时出现错误',
      });
      setLoading(false);
    }
  };

  const handleConfirm = async (confirmed: boolean) => {
    setConfirmModalVisible(false);
    if (!confirmed || !currentSession) return;

    setLoading(true);
    try {
      const res = await confirmAction({
        session_id: currentSession.id,
        confirm: true,
      });
      setCurrentSession(res.data.session);
      setSessions((prev) =>
        prev.map((s) => (s.id === res.data.session.id ? res.data.session : s))
      );
    } catch (error: any) {
      console.error('确认操作失败', error);
      Modal.error({
        title: '操作失败',
        content: error.msg || '执行操作时出现错误',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleQuickCommand = (prompt: string) => {
    if (!currentSession) {
      handleCreateSession().then(() => {
        setInputValue(prompt);
      });
    } else {
      setInputValue(prompt);
    }
  };

  const formatTime = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));

    if (days === 0) {
      return '今天';
    } else if (days === 1) {
      return '昨天';
    } else if (days < 7) {
      return `${days}天前`;
    } else {
      return date.toLocaleDateString('zh-CN', {
        month: 'short',
        day: 'numeric',
      });
    }
  };

  const renderMessageContent = (message: ChatMessage) => {
    const content = message?.content || '';
    
    if (!content || typeof content !== 'string') {
      return <div className="message-content"><p>-</p></div>;
    }

    if (content.length > 50000) {
      return <div className="message-content"><p>{content.substring(0, 50000)}...</p></div>;
    }

    return (
      <div className="message-content">
        <ReactMarkdown
          remarkPlugins={[remarkGfm]}
          rehypePlugins={[rehypeHighlight]}
          components={{
            p: ({ children }) => <p className="markdown-p">{children}</p>,
            code: ({ className, children }) => {
              const lang = className?.replace(/language-/, '') || '';
              return (
                <code className={`language-${lang} ${className || ''}`}>
                  {children}
                </code>
              );
            },
            pre: ({ children }) => (
              <pre className="markdown-pre"><code>{children}</code></pre>
            ),
            blockquote: ({ children }) => (
              <blockquote className="markdown-quote">{children}</blockquote>
            ),
            ul: ({ children }) => <ul className="markdown-ul">{children}</ul>,
            ol: ({ children }) => <ol className="markdown-ol">{children}</ol>,
            li: ({ children }) => <li className="markdown-li">{children}</li>,
            strong: ({ children }) => <strong className="markdown-strong">{children}</strong>,
            table: ({ children }) => (
              <div className="markdown-table-container">
                <table className="markdown-table">{children}</table>
              </div>
            ),
            th: ({ children }) => <th className="markdown-th">{children}</th>,
            td: ({ children }) => <td className="markdown-td">{children}</td>,
          }}
        >
          {content}
        </ReactMarkdown>
      </div>
    );
  };

  return (
    <Layout className="assistant-layout">
      <Sider className="assistant-sider" width={280}>
        <div className="sider-header">
          <div className="sider-title">
            <MessageOutlined />
            <span>开发助手</span>
          </div>
          <Button
            type="primary"
            size="small"
            icon={<PlusOutlined />}
            onClick={handleCreateSession}
          >
            新建会话
          </Button>
        </div>

        <List
          dataSource={sessions}
          renderItem={(session) => (
            <List.Item
              className={`session-item ${currentSession?.id === session.id ? 'active' : ''}`}
              onClick={() => handleSelectSession(session.id)}
              actions={[
                <Tooltip title="删除">
                  <Button
                    type="text"
                    danger
                    icon={<DeleteOutlined />}
                    onClick={(e) => {
                      e.stopPropagation();
                      handleDeleteSession(session.id);
                    }}
                  />
                </Tooltip>,
              ]}
            >
              <div className="session-info">
                <div className="session-title">{session.title || '未命名会话'}</div>
                <div className="session-meta">
                  <span className="session-user">
                    <UserOutlined style={{ fontSize: 10 }} />
                    {session.user_name || '未知用户'}
                  </span>
                  <span className="session-time">
                    {formatDate(session.updated_at || session.created_at)}
                  </span>
                </div>
              </div>
            </List.Item>
          )}
          locale={{ emptyText: <Empty description="暂无会话" /> }}
        />
      </Sider>

      <Content className="assistant-content">
        {currentSession ? (
          <div className="chat-container">
            <div ref={scrollRef} className="chat-messages">
              {currentSession?.messages?.map((message) => (
                <div
                  key={message.id}
                  className={`message-item ${message.role}`}
                >
                  <div className="message-avatar">
                    {message.role === 'user' ? (
                      <UserOutlined />
                    ) : (
                      <RobotOutlined />
                    )}
                  </div>
                  <div className="message-body">
                    {message.thought && (
                      <div className="message-thought">
                        <BulbOutlined className="thought-icon" />
                        <span className="thought-text">{message.thought}</span>
                      </div>
                    )}
                    {renderMessageContent(message)}
                    {message.tool_calls && message.tool_calls.length > 0 && (
                      <div className="tool-calls">
                        <div className="tool-calls-title">工具调用:</div>
                        {message.tool_calls.map((tool, idx) => (
                          <div key={idx} className="tool-call-item">
                            <Tag color="blue">{tool.tool_name}</Tag>
                            <code className="tool-params">
                              {JSON.stringify(tool.parameters, null, 2)}
                            </code>
                          </div>
                        ))}
                      </div>
                    )}
                    {message.tool_results && message.tool_results.length > 0 && (
                      <div className="tool-results">
                        <div className="tool-results-title">执行结果:</div>
                        {message.tool_results.map((result, idx) => (
                          <div key={idx} className={`tool-result-item ${result.success ? 'success' : 'error'}`}>
                            <Tag color={result.success ? 'green' : 'red'}>{result.tool_name}</Tag>
                            <pre className="tool-result-output">{result.success ? result.output : result.error}</pre>
                          </div>
                        ))}
                      </div>
                    )}
                    <div className="message-time">{formatTime(message.created_at)}</div>
                  </div>
                </div>
              ))}

              {loading && (
                <div className="loading-message">
                  <RobotOutlined className="bot-icon" />
                  <div className="loading-dots">
                    <Spin size="small" />
                    <span>正在思考...</span>
                  </div>
                </div>
              )}
              <div ref={setScrollBottom} className="scroll-bottom" />
            </div>

            <div className="chat-input-area">
              <div className="quick-commands">
                {quickCommands.map((cmd) => (
                  <Button
                    key={cmd.label}
                    size="small"
                    ghost
                    type="primary"
                    onClick={() => handleQuickCommand(cmd.prompt)}
                  >
                    {cmd.icon}
                    {cmd.label}
                  </Button>
                ))}
              </div>

              <div className="input-wrapper">
                <TextArea
                  value={inputValue}
                  onChange={(e) => setInputValue(e.target.value)}
                  onPressEnter={(e) => {
                    if (e.ctrlKey || e.metaKey) {
                      e.preventDefault();
                      handleSendMessage();
                    }
                  }}
                  placeholder="输入你的问题... (Ctrl+Enter 发送)"
                  rows={3}
                  disabled={loading}
                />
                <Button
                  type="primary"
                  icon={<SendOutlined />}
                  onClick={handleSendMessage}
                  disabled={loading || !inputValue.trim()}
                  className="send-btn"
                >
                  发送
                </Button>
              </div>
            </div>
          </div>
        ) : (
          <div className="empty-state">
            <Card className="empty-card">
              <div className="empty-icon">
                <MessageOutlined />
              </div>
              <h3>欢迎使用开发助手</h3>
              <p>我是您的全栈开发助手，已加载 Gin/UniApp/Vue3 开发环境。</p>
              <p>我可以帮您编写代码、调试接口、生成图表。</p>
              <Button type="primary" size="large" onClick={handleCreateSession}>
                开始对话
              </Button>
            </Card>
          </div>
        )}
      </Content>

      <Modal
        title="确认执行操作"
        open={confirmModalVisible}
        onCancel={() => setConfirmModalVisible(false)}
        footer={[
          <Button key="back" onClick={() => handleConfirm(false)}>
            取消
          </Button>,
          <Button key="submit" type="primary" onClick={() => handleConfirm(true)}>
            确认执行
          </Button>,
        ]}
      >
        <p>此操作涉及以下文件的修改：</p>
        <List
          dataSource={currentSession?.affected_files || []}
          renderItem={(file) => (
            <List.Item>
              <FileTextOutlined />
              <span>{file}</span>
            </List.Item>
          )}
          locale={{ emptyText: '暂无影响文件' }}
        />
        <p className="confirm-warning">
          请仔细检查以上文件，确认无误后点击"确认执行"。此操作无法撤销。
        </p>
      </Modal>

      <Modal
        title="新建会话"
        open={createSessionModalVisible}
        onCancel={() => setCreateSessionModalVisible(false)}
        footer={[
          <Button key="back" onClick={() => setCreateSessionModalVisible(false)}>
            取消
          </Button>,
          <Button key="submit" type="primary" onClick={handleConfirmCreateSession}>
            创建
          </Button>,
        ]}
      >
        <p>请输入会话标题：</p>
        <Input
          value={newSessionTitle}
          onChange={(e) => setNewSessionTitle(e.target.value)}
          placeholder="例如：订单模块开发"
          maxLength={50}
        />
      </Modal>
    </Layout>
  );
};
