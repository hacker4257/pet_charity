import { useState } from 'react';
import { Card, Tabs, Form, Input, Button, message } from 'antd';
import { UserOutlined, LockOutlined, MobileOutlined } from '@ant-design/icons';
import { useNavigate, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { login, sendSmsCode, smsLogin } from '../../api/auth';
import useAuthStore from '../../store/useAuthStore';

const Login = () => {
    const navigate = useNavigate();
    const { t }  = useTranslation();
    const { loginSuccess } = useAuthStore();

    //账号密码登录
    const [accountLoading, setAccountLoading] = useState(false)

    const onAccountLogin = async (values: {account: string; password: string}) => {
        setAccountLoading(true);
        try {
            const res: any = await login(values);
            await loginSuccess(res.data.access_token, res.data.refresh_token);
            message.success('Login successful');
            navigate('/');
        }catch{
            //
        }finally {
            setAccountLoading(false)
        }
    };

    //验证码
    const [smsLoading, setSmsLoading] = useState(false);
    const [countdown, setCountdown] = useState(0);

    //发送验证码
    const handleSendCode = async (phone:string) => {
        if (!phone || phone.length != 11) {
            message.warning('Please enter a valid phone number');
            return;
        }

        try {
            await sendSmsCode({phone, purpose:'login'});
            message.success('Code sent');

            //60seconds
            setCountdown(60);
            const timer = setInterval(()=> {
                setCountdown((prev)=> {
                    if (prev <= 1) {
                        clearInterval(timer);
                        return 0;
                    }
                    return prev - 1
                });
            }, 1000)
        }catch {
            //
        }
    };

    const onSmsLogin = async (values:{phone: string; code: string}) => {
        setSmsLoading(true);
        try{
            const res: any = await smsLogin(values);
            await loginSuccess(res.data.access_token, res.data.refresh_token);
            message.success('Login successful');
            navigate('/');
        }catch{
            //
        }finally{
            setSmsLoading(false);
        }
    };

    const [smsForm] = Form.useForm();

    return (
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center',minHeight: '100vh', background: '#f0f2f5' }}>
            <Card style={{ width: 420 }}>
            <h2 style={{ textAlign: 'center', marginBottom: 24 }}>Pet Charity</h2>

            <Tabs
                centered
                items={[
                {
                    key: 'account',
                    label: t('common.login'),
                    children: (
                    <Form onFinish={onAccountLogin} size="large">
                        <Form.Item name="account" rules={[{ required: true, message:'Please enter username or email' }]}>
                        <Input prefix={<UserOutlined />} placeholder="Username or Email"/>
                        </Form.Item>
                        <Form.Item name="password" rules={[{ required: true, message:'Please enter password' }]}>
                        <Input.Password prefix={<LockOutlined />} placeholder="Password"/>
                        </Form.Item>

                        <Form.Item>
                        <Button type="primary" htmlType="submit"loading={accountLoading} block>
                            {t('common.login')}
                        </Button>
                        </Form.Item>
                    </Form>
                    ),
                },
                {
                    key: 'sms',
                    label: 'SMS Login',
                    children: (
                    <Form form={smsForm} onFinish={onSmsLogin} size="large">
                        <Form.Item name="phone" rules={[{ required: true, len: 11,message: 'Please enter 11-digit phone number' }]}>
                        <Input prefix={<MobileOutlined />} placeholder="Phone Number" maxLength={11} />
                        </Form.Item>

                        <Form.Item name="code" rules={[{ required: true, len: 6, message:'Please enter 6-digit code' }]}>
                        <div style={{ display: 'flex', gap: 8 }}>
                            <Input placeholder="Verification Code" maxLength={6} />
                            <Button
                            disabled={countdown > 0}
                            onClick={() =>handleSendCode(smsForm.getFieldValue('phone'))}
                            style={{ width: 120 }}
                            >
                            {countdown > 0 ? `${countdown}s` : 'Send Code'}
                            </Button>
                        </div>
                        </Form.Item>

                        <Form.Item>
                        <Button type="primary" htmlType="submit" loading={smsLoading} block>
                            {t('common.login')}
                        </Button>
                        </Form.Item>
                    </Form>
                    ),
                },
                ]}
            />

            <div style={{ textAlign: 'center' }}>
                No account? <Link to="/register">{t('common.register')}</Link>
            </div>
            </Card>
        </div>
        );
}

export default Login