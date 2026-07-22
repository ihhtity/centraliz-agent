import { Form, Input, Button, Card, message } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { login } from '@/api';
import './Login.scss';

export const Login = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();

  const handleSubmit = async (values: { username: string; password: string; type: string }) => {
    try {
      const res = await login({
        account: values.username,
        password: values.password,
        type: 'account'
      });
      localStorage.setItem('token', res.data.token);
      message.success('登录成功');
      navigate('/dashboard');
    } catch (err: any) {
      message.error(err.msg || '登录失败，请检查用户名和密码');
    }
  };

  return (
    <div className="login-container">
      <Card className="login-card">
        <h2 className="login-title">管理后台</h2>
        <Form form={form} onFinish={handleSubmit} layout="vertical">
          <Form.Item
            name="username"
            label="用户名"
            initialValue="17727293262"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input prefix={<UserOutlined />} placeholder="请输入用户名" />
          </Form.Item>
          <Form.Item
            name="password"
            label="密码"
            initialValue="12345678"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="请输入密码" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              登录
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};
