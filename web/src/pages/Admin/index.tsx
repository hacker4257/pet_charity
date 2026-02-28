import { useEffect, useState } from 'react';
import { Tabs, Card, Row, Col, Statistic, Table, Tag, Button, Modal, Input, message} from 'antd';
import {
    DashboardOutlined, BankOutlined, UserOutlined,
    TeamOutlined, HeartOutlined, DollarOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { getAdminStats, getPendingOrgs, reviewOrg, getUserList, updateUserRole, updateUserStatus } from '../../api/admin';

const Admin = () => {
    const { t } = useTranslation();

    // ========== Tab 1: 数据总览 ==========
    const [stats, setStats] = useState<any>(null);

    const fetchStats = async () => {
        try {
            const res: any = await getAdminStats();
            setStats(res.data);
        } catch {
            //
        }
    };

    useEffect(() => {
        fetchStats();
    }, []);

    // ========== Tab 2: 机构审核 ==========
    const [orgs, setOrgs] = useState<any[]>([]);
    const [orgTotal, setOrgTotal] = useState(0);
    const [orgPage, setOrgPage] = useState(1);
    const [orgLoading, setOrgLoading] = useState(false);

    const fetchOrgs = async () => {
        setOrgLoading(true);
        try {
            const res: any = await getPendingOrgs({ page: orgPage, page_size: 10 });
            setOrgs(res.data.list || []);
            setOrgTotal(res.data.total || 0);
        } catch {
            //
        } finally {
            setOrgLoading(false);
        }
    };

    useEffect(() => {
        fetchOrgs();
    }, [orgPage]);

    // 通过机构
    const handleApprove = (id: number, name: string) => {
        Modal.confirm({
            title: '确认通过',
            content: `确定通过「${name}」的机构注册申请吗？`,
            okText: '通过',
            onOk: async () => {
                await reviewOrg(id, { status: 'approved' });
                message.success('已通过');
                fetchOrgs();
            },
        });
    };

    // 拒绝机构
    const handleReject = (id: number, name: string) => {
        let reason = '';
        Modal.confirm({
            title: '拒绝原因',
            content: (
                <Input.TextArea
                    rows={3}
                    placeholder="请输入拒绝原因"
                    onChange={(e) => { reason = e.target.value; }}
                />
            ),
            okText: '拒绝',
            okType: 'danger',
            onOk: async () => {
                if (!reason.trim()) {
                    message.warning('请输入拒绝原因');
                    throw new Error('need reason'); // 阻止弹窗关闭
                }
                await reviewOrg(id, { status: 'rejected', reject_reason: reason });
                message.success('已拒绝');
                fetchOrgs();
            },
        });
    };

    const orgColumns = [
        { title: '机构名称', dataIndex: 'name', key: 'name' },
        { title: '联系电话', dataIndex: 'contact_phone', key: 'phone' },
        { title: '地址', dataIndex: 'address', key: 'address', ellipsis: true },
        { title: '容纳量', dataIndex: 'capacity', key: 'capacity' },
        {
            title: '申请时间',
            dataIndex: 'created_at',
            key: 'createdAt',
            render: (t: string) => new Date(t).toLocaleDateString(),
        },
        {
            title: '操作',
            key: 'action',
            render: (_: any, record: any) => (
                <div style={{ display: 'flex', gap: 8 }}>
                    <Button type="primary" size="small" onClick={() =>handleApprove(record.id, record.name)}>
                        通过
                    </Button>
                    <Button danger size="small" onClick={() =>handleReject(record.id, record.name)}>
                        拒绝
                    </Button>
                </div>
            ),
        },
    ];

    // ========== Tab 3: 用户管理 ==========
    const [users, setUsers] = useState<any[]>([]);
    const [userTotal, setUserTotal] = useState(0);
    const [userPage, setUserPage] = useState(1);
    const [userLoading, setUserLoading] = useState(false);

    const fetchUsers = async () => {
        setUserLoading(true);
        try {
            const res: any = await getUserList({ page: userPage, page_size: 10 });
            setUsers(res.data.list || []);
            setUserTotal(res.data.total || 0);
        } catch {
            //
        } finally {
            setUserLoading(false);
        }
    };

    useEffect(() => {
        fetchUsers();
    }, [userPage]);

    const roleColorMap: Record<string, string> = {
        admin: 'red',
        org: 'blue',
        user: 'default',
    };

    // 切换用户状态
    const toggleUserStatus = async (record: any) => {
        const newStatus = record.status === 'active' ? 'disabled' : 'active';
        const label = newStatus === 'disabled' ? '禁用' : '启用';

        Modal.confirm({
            title: `确认${label}`,
            content: `确定${label}用户「${record.username}」吗？`,
            okType: newStatus === 'disabled' ? 'danger' : 'primary',
            onOk: async () => {
                await updateUserStatus(record.id, { status: newStatus });
                message.success(`已${label}`);
                fetchUsers();
            },
        });
    };

    // 修改角色
    const handleChangeRole = (record: any) => {
        let newRole = record.role;
        Modal.confirm({
            title: '修改角色',
            content: (
                <div style={{ marginTop: 12 }}>
                    <span>当前角色：{record.role}，修改为：</span>
                    <select
                        defaultValue={record.role}
                        onChange={(e) => { newRole = e.target.value; }}
                        style={{ marginLeft: 8, padding: '4px 8px' }}
                    >
                        <option value="user">user</option>
                        <option value="org">org</option>
                        <option value="admin">admin</option>
                    </select>
                </div>
            ),
            onOk: async () => {
                if (newRole === record.role) return;
                await updateUserRole(record.id, { role: newRole });
                message.success('角色已修改');
                fetchUsers();
            },
        });
    };

    const userColumns = [
        { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
        { title: '用户名', dataIndex: 'username', key: 'username' },
        { title: '昵称', dataIndex: 'nickname', key: 'nickname' },
        { title: '邮箱', dataIndex: 'email', key: 'email', ellipsis: true },
        {
            title: '角色',
            dataIndex: 'role',
            key: 'role',
            render: (role: string) => <Tag color={roleColorMap[role]}>{role}</Tag>,
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            render: (s: string) => (
                <Tag color={s === 'active' ? 'green' : 'red'}>
                    {s === 'active' ? '正常' : '禁用'}
                </Tag>
            ),
        },
        {
            title: '注册时间',
            dataIndex: 'created_at',
            key: 'createdAt',
            render: (t: string) => new Date(t).toLocaleDateString(),
        },
        {
            title: '操作',
            key: 'action',
            render: (_: any, record: any) => (
                <div style={{ display: 'flex', gap: 8 }}>
                    <Button size="small" onClick={() => handleChangeRole(record)}>
                        改角色
                    </Button>
                    <Button
                        size="small"
                        danger={record.status === 'active'}
                        onClick={() => toggleUserStatus(record)}
                    >
                        {record.status === 'active' ? '禁用' : '启用'}
                    </Button>
                </div>
            ),
        },
    ];

    return (
        <div>
            <h2>管理后台</h2>

            <Tabs items={[
                {
                    key: 'dashboard',
                    label: <span><DashboardOutlined /> 数据总览</span>,
                    children: (
                        <Row gutter={24}>
                            <Col span={6}>
                                <Card>
                                    <Statistic
                                        title="注册用户"
                                        value={stats?.user_count || 0}
                                        prefix={<TeamOutlined />}
                                    />
                                </Card>
                            </Col>
                            <Col span={6}>
                                <Card>
                                    <Statistic
                                        title="认证机构"
                                        value={stats?.org_count || 0}
                                        prefix={<BankOutlined />}
                                    />
                                </Card>
                            </Col>
                            <Col span={6}>
                                <Card>
                                    <Statistic
                                        title="待领养宠物"
                                        value={stats?.pet_stats?.adoptable_pets ||0}
                                        prefix={<HeartOutlined />}
                                    />
                                </Card>
                            </Col>
                            <Col span={6}>
                                <Card>
                                    <Statistic
                                        title="捐赠总额（元）"
                                        value={((stats?.donation_stats?.total_amount|| 0) / 100).toFixed(2)}
                                        prefix={<DollarOutlined />}
                                    />
                                </Card>
                            </Col>
                        </Row>
                    ),
                },
                {
                    key: 'orgs',
                    label: <span><BankOutlined /> 机构审核</span>,
                    children: (
                        <Card>
                            <Table
                                dataSource={orgs}
                                columns={orgColumns}
                                rowKey="id"
                                loading={orgLoading}
                                pagination={{
                                    current: orgPage,
                                    pageSize: 10,
                                    total: orgTotal,
                                    onChange: (p) => setOrgPage(p),
                                }}
                            />
                        </Card>
                    ),
                },
                {
                    key: 'users',
                    label: <span><UserOutlined /> 用户管理</span>,
                    children: (
                        <Card>
                            <Table
                                dataSource={users}
                                columns={userColumns}
                                rowKey="id"
                                loading={userLoading}
                                pagination={{
                                    current: userPage,
                                    pageSize: 10,
                                    total: userTotal,
                                    onChange: (p) => setUserPage(p),
                                }}
                            />
                        </Card>
                    ),
                },
            ]} />
        </div>
    );
};

export default Admin;