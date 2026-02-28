import { useState } from 'react';
import { Card, Form, Input, Button, message } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons';
import { useNavigate, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { register } from '../../api/auth';

const Register = () => {
const navigate = useNavigate();
const { t } = useTranslation();
const [loading, setLoading] = useState(false);

const onFinish = async (values: {
    username: string;
    email: string;
    password: string;
    confirm: string;
    nickname: string;
    }) => {
        setLoading(true);
        try {
            await register({
                username: values.username,
                email: values.email,
                password: values.password,
            });
            message.success('Registration successful, please login');
            navigate('/login');
        } catch {
        // 错误已拦截
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center',minHeight: '100vh', background: '#f0f2f5' }}>
        <Card style={{ width: 420 }}>
            <h2 style={{ textAlign: 'center', marginBottom: 24 }}>Create Account</h2>

            <Form onFinish={onFinish} size="large">
            <Form.Item
                name="username"
                rules={[
                { required: true, message: 'Please enter username' },
                { min: 3, max: 50, message: '3-50 characters' },
                ]}
            >
                <Input prefix={<UserOutlined />} placeholder="Username" />
            </Form.Item>

            <Form.Item
                name="email"
                rules={[
                { required: true, message: 'Please enter email' },
                { type: 'email', message: 'Invalid email format' },
                ]}
            >
                <Input prefix={<MailOutlined />} placeholder="Email" />
            </Form.Item>

            <Form.Item name="nickname">
                <Input prefix={<UserOutlined />} placeholder="Nickname (optional)" />
            </Form.Item>

            <Form.Item
                name="password"
                rules={[
                { required: true, message: 'Please enter password' },
                { min: 6, max: 50, message: '6-50 characters' },
                ]}
            >
                <Input.Password prefix={<LockOutlined />} placeholder="Password" />
            </Form.Item>

            <Form.Item
                name="confirm"
                dependencies={['password']}
                rules={[
                { required: true, message: 'Please confirm password' },
                ({ getFieldValue }) => ({
                    validator(_, value) {
                    if (!value || getFieldValue('password') === value) {
                        return Promise.resolve();
                    }
                    return Promise.reject(new Error('Passwords do not match'));
                    },
                }),
                ]}
            >
                <Input.Password prefix={<LockOutlined />} placeholder="Confirm Password"/>
            </Form.Item>

            <Form.Item>
                <Button type="primary" htmlType="submit" loading={loading} block>
                {t('common.register')}
                </Button>
            </Form.Item>
            </Form>

            <div style={{ textAlign: 'center' }}>
                Already have an account? <Link to="/login">{t('common.login')}</Link>
            </div>
        </Card>
        </div>
    );
};

export default Register;
