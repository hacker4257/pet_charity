import { useEffect, useState } from 'react';
import { Tabs, Card, Form, Input, Button, Select, Upload, Avatar, Table, Tag,message } from 'antd';
import { UserOutlined, LockOutlined, HeartOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import useAuthStore from '../../store/useAuthStore';
import { updateProfile, changePassword, uploadeAvatar } from '../../api/user';
import { getMyAdoptions } from '../../api/adoption';
import { BankOutlined } from '@ant-design/icons';
import OrgDashboard from '../../components/OrgDashboard';


const Profile = () => {
    const { t, i18n } = useTranslation();
    const { user, fetchUser, logout } = useAuthStore();
    const [profileForm] = Form.useForm();
    const [passwordForm] = Form.useForm();

    // ========== Tab 1: 基本资料 ==========
    const [saving, setSaving] = useState(false);

    // 用户数据加载后填入表单
    useEffect(() => {
        if (user) {
            profileForm.setFieldsValue({
                nickname: user.nickname,
                email: user.email,
                phone: user.phone,
                language: user.language || 'zh-CN',
            });
        }
    }, [user]);

    const onSaveProfile = async (values: any) => {
        setSaving(true);
        try {
            await updateProfile(values);
            message.success(t('profile.saveSuccess'));

            // 如果切换了语言，同步切换前端
            if (values.language && values.language !== i18n.language) {
                i18n.changeLanguage(values.language);
            }

            // 重新拉取用户信息
            await fetchUser();
        } catch {
            //
        } finally {
            setSaving(false);
        }
    };

    // 头像上传
    const handleAvatarUpload = async (file: File) => {
        try {
            await uploadeAvatar(file);
            message.success('头像已更新');
            await fetchUser();
        } catch {
            //
        }
        return false; // 阻止 antd 默认上传行为
    };

    // ========== Tab 2: 修改密码 ==========
    const [changingPwd, setChangingPwd] = useState(false);

    const onChangePassword = async (values: any) => {
        setChangingPwd(true);
        try {
            await changePassword({ old_password: values.old_password, new_password: values.new_password });
            message.success(t('profile.passwordChanged'));
            passwordForm.resetFields();

            // 修改密码后自动登出
            setTimeout(() => logout(), 1500);
        } catch {
            //
        } finally {
            setChangingPwd(false);
        }
    };

    // ========== Tab 3: 我的领养 ==========
    const [adoptions, setAdoptions] = useState<any[]>([]);
    const [adoptTotal, setAdoptTotal] = useState(0);
    const [adoptPage, setAdoptPage] = useState(1);
    const [adoptLoading, setAdoptLoading] = useState(false);

    const fetchAdoptions = async () => {
        setAdoptLoading(true);
        try {
            const res: any = await getMyAdoptions({ page: adoptPage, page_size: 10
});
            setAdoptions(res.data.list || []);
            setAdoptTotal(res.data.total || 0);
        } catch {
            //
        } finally {
            setAdoptLoading(false);
        }
    };

    useEffect(() => {
        fetchAdoptions();
    }, [adoptPage]);

    // 领养状态映射
    const adoptionStatusMap: Record<string, { color: string; text: string }> = {
        pending: { color: 'orange', text: '审核中' },
        approved: { color: 'green', text: '已通过' },
        rejected: { color: 'red', text: '已拒绝' },
        completed: { color: 'blue', text: '已完成' },
    };

    const adoptionColumns = [
        {
            title: '宠物',
            dataIndex: ['pet', 'name'],
            key: 'petName',
            render: (name: string) => name || '-',
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            render: (s: string) => {
                const item = adoptionStatusMap[s] || { color: 'default', text: s };
                return <Tag color={item.color}>{item.text}</Tag>;
            },
        },
        {
            title: '申请时间',
            dataIndex: 'created_at',
            key: 'createdAt',
            render: (t: string) => new Date(t).toLocaleDateString(),
        },
        {
            title: '拒绝原因',
            dataIndex: 'reject_reason',
            key: 'rejectReason',
            render: (r: string) => r || '-',
        },
    ];

    return (
        <div style={{ maxWidth: 800, margin: '0 auto' }}>
            <h2>{t('profile.title')}</h2>

            <Tabs items={[
                {
                    key: 'basic',
                    label: <span><UserOutlined /> {t('profile.basicInfo')}</span>,
                    children: (
                        <Card>
                            {/* 头像 */}
                            <div style={{ textAlign: 'center', marginBottom: 24 }}>
                                <Upload
                                    showUploadList={false}
                                    beforeUpload={handleAvatarUpload}
                                    accept="image/*"
                                >
                                    <Avatar
                                        size={80}
                                        src={user?.avatar}
                                        icon={<UserOutlined />}
                                        style={{ cursor: 'pointer' }}
                                    />
                                    <div style={{ marginTop: 8, color: '#1890ff',cursor: 'pointer' }}>
                                        {t('profile.clickUpload')}
                                    </div>
                                </Upload>
                            </div>

                            {/* 资料表单 */}
                            <Form
                                form={profileForm}
                                layout="vertical"
                                onFinish={onSaveProfile}
                            >
                                <Form.Item label={t('profile.nickname')} name="nickname">
                                    <Input placeholder={t('profile.nickname')} />
                                </Form.Item>

                                <Form.Item
                                    label={t('profile.email')}
                                    name="email"
                                    rules={[{ type: 'email', message: '邮箱格式不正确' }]}
                                >
                                    <Input placeholder={t('profile.email')} />
                                </Form.Item>

                                <Form.Item label={t('profile.phone')} name="phone">
                                    <Input placeholder={t('profile.phone')} maxLength={11} />
                                </Form.Item>

                                <Form.Item label={t('profile.language')} name="language">
                                    <Select>
                                        <Select.Option value="zh-CN">中文</Select.Option>
                                        <Select.Option value="en-US">English</Select.Option>
                                    </Select>
                                </Form.Item>

                                <Form.Item>
                                    <Button type="primary" htmlType="submit" loading={saving} block>
                                        {t('common.save')}
                                    </Button>
                                </Form.Item>
                            </Form>
                        </Card>
                    ),
                },
                {
                    key: 'password',
                    label: <span><LockOutlined /> {t('profile.changePassword')}</span>,
                    children: (
                        <Card>
                            <Form
                                form={passwordForm}
                                layout="vertical"
                                onFinish={onChangePassword}
                            >
                                <Form.Item
                                    name="old_password"
                                    label={t('profile.oldPassword')}
                                    rules={[
                                        { required: true, message: '请输入当前密码' },
                                    ]}
                                >
                                    <Input.Password prefix={<LockOutlined />} />
                                </Form.Item>

                                <Form.Item
                                    name="new_password"
                                    label={t('profile.newPassword')}
                                    rules={[
                                        { required: true, message: '请输入新密码' },
                                        { min: 6, max: 50, message: '6-50 个字符' },
                                    ]}
                                >
                                    <Input.Password prefix={<LockOutlined />} />
                                </Form.Item>

                                <Form.Item
                                    name="confirm"
                                    label={t('profile.confirmPassword')}
                                    dependencies={['new_password']}
                                    rules={[
                                        { required: true, message: '请确认新密码' },
                                        ({ getFieldValue }) => ({
                                            validator(_, value) {
                                                if (!value || getFieldValue('new_password') === value) {
                                                    return Promise.resolve();
                                                }
                                                return Promise.reject(new Error('两次密码不一致'));
                                            },
                                        }),
                                    ]}
                                >
                                    <Input.Password prefix={<LockOutlined />} />
                                </Form.Item>

                                <Form.Item>
                                    <Button type="primary" htmlType="submit" loading={changingPwd} block>
                                        {t('profile.changePassword')}
                                    </Button>
                                </Form.Item>
                            </Form>
                        </Card>
                    ),
                },
                {
                    key: 'adoptions',
                    label: <span><HeartOutlined />{t('profile.myAdoptions')}</span>,
                    children: (
                        <Card>
                            <Table
                                dataSource={adoptions}
                                columns={adoptionColumns}
                                rowKey="id"
                                loading={adoptLoading}
                                pagination={{
                                    current: adoptPage,
                                    pageSize: 10,
                                    total: adoptTotal,
                                    onChange: (p) => setAdoptPage(p),
                                }}
                            />
                        </Card>
                    ),
                },
                {
                    key: 'org',
                    label: <span><BankOutlined /> {t('profile.orgManage')}</span>,
                    children: <OrgDashboard />,
                },
            ]} />
        </div>
    );
};

export default Profile;