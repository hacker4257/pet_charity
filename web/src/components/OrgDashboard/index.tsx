import { useEffect, useState } from 'react';
import { Card, Tabs, Form, Input, InputNumber, Button, Table, Tag, Modal, Select, message } from 'antd';
import { PlusOutlined, UploadOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { getMyOrg, createOrg, updateOrg } from '../../api/organization';
import { getPets, createPet, updatePet, deletePet, uploadPetImage } from '../../api/pet';
import { getOrgAdoptions, reviewAdoption, completeAdoption } from '../../api/adoption';

// ========== 机构注册表单 ==========
const OrgRegisterForm = ({ onSuccess }: { onSuccess: () => void }) => {
    const [loading, setLoading] = useState(false);

    const onFinish = async (values: any) => {
        setLoading(true);
        try {
            await createOrg(values);
            message.success('注册成功，等待管理员审核');
            onSuccess();
        } catch {
            //
        } finally {
            setLoading(false);
        }
    };

    return (
        <Card title="注册救助机构">
            <Form layout="vertical" onFinish={onFinish} style={{ maxWidth: 500 }}>
                <Form.Item name="name" label="机构名称" rules={[{ required: true}]}>
                    <Input placeholder="如：阳光流浪动物救助站" />
                </Form.Item>
                <Form.Item name="description" label="机构简介" rules={[{ required:true }]}>
                    <Input.TextArea rows={3} placeholder="介绍机构的宗旨、主要业务等" />
                </Form.Item>
                <Form.Item name="address" label="地址" rules={[{ required: true }]}>
                    <Input placeholder="详细地址" />
                </Form.Item>
                <Form.Item name="contact_phone" label="联系电话" rules={[{ required:true }]}>
                    <Input placeholder="联系电话" maxLength={11} />
                </Form.Item>
                <Form.Item name="capacity" label="最大容纳量">
                    <InputNumber min={0} placeholder="可容纳动物数量" style={{ width: '100%' }} />
                </Form.Item>
                <Form.Item>
                    <Button type="primary" htmlType="submit" loading={loading} block>
                        提交注册
                    </Button>
                </Form.Item>
            </Form>
        </Card>
    );
};

// ========== 机构信息编辑 ==========
const OrgInfoTab = ({ org, onUpdate }: { org: any; onUpdate: () => void }) => {
    const [form] = Form.useForm();
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        form.setFieldsValue(org);
    }, [org]);

    const onFinish = async (values: any) => {
        setSaving(true);
        try {
            await updateOrg(org.id, values);
            message.success('已更新');
            onUpdate();
        } catch {
            //
        } finally {
            setSaving(false);
        }
    };

    return (
        <Form form={form} layout="vertical" onFinish={onFinish} style={{ maxWidth:500 }}>
            <Form.Item label="机构名称">
                <Input value={org.name} disabled />
            </Form.Item>
            <Form.Item name="description" label="简介">
                <Input.TextArea rows={3} />
            </Form.Item>
            <Form.Item name="address" label="地址">
                <Input />
            </Form.Item>
            <Form.Item name="contact_phone" label="联系电话">
                <Input maxLength={11} />
            </Form.Item>
            <Form.Item name="capacity" label="容纳量">
                <InputNumber min={0} style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item>
                <Button type="primary" htmlType="submit" loading={saving}>保存</Button>
            </Form.Item>
        </Form>
    );
};

// ========== 宠物管理 ==========
const PetManageTab = ({ orgId }: { orgId: number }) => {
    const { t } = useTranslation();
    const [pets, setPets] = useState<any[]>([]);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);
    const [loading, setLoading] = useState(false);
    const [modalOpen, setModalOpen] = useState(false);
    const [editingPet, setEditingPet] = useState<any>(null);
    const [form] = Form.useForm();
    const [submitting, setSubmitting] = useState(false);

    const fetchPets = async () => {
        setLoading(true);
        try {
            const res: any = await getPets({ org_id: orgId, page, pageSize: 10 });
            setPets(res.data.list || []);
            setTotal(res.data.total || 0);
        } catch {
            //
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => { fetchPets(); }, [page]);

    // 打开新增/编辑弹窗
    const openModal = (pet?: any) => {
        setEditingPet(pet || null);
        if (pet) {
            form.setFieldsValue(pet);
        } else {
            form.resetFields();
        }
        setModalOpen(true);
    };

    // 提交表单
    const onFinish = async (values: any) => {
        setSubmitting(true);
        try {
            if (editingPet) {
                await updatePet(editingPet.id, values);
                message.success('已更新');
            } else {
                await createPet(values);
                message.success('已添加');
            }
            setModalOpen(false);
            fetchPets();
        } catch {
            //
        } finally {
            setSubmitting(false);
        }
    };

    // 删除
    const handleDelete = (id: number, name: string) => {
        Modal.confirm({
            title: '确认删除',
            content: `确定删除「${name}」吗？此操作不可恢复。`,
            okType: 'danger',
            onOk: async () => {
                await deletePet(id);
                message.success('已删除');
                fetchPets();
            },
        });
    };

    // 上传图片
    const handleUploadImage = (petId: number) => {
        // 创建隐藏的 input
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = 'image/*';
        input.onchange = async (e: any) => {
            const file = e.target.files?.[0];
            if (!file) return;
            try {
                await uploadPetImage(petId, file);
                message.success('图片已上传');
            } catch {
                //
            }
        };
        input.click();
    };

    const statusMap: Record<string, { color: string; text: string }> = {
        adoptable: { color: 'green', text: t('pet.available') },
        reserved: { color: 'orange', text: t('pet.reserved') },
        adopted: { color: 'blue', text: t('pet.adopted') },
    };

    const columns = [
        { title: '名字', dataIndex: 'name', key: 'name' },
        { title: '物种', dataIndex: 'species', key: 'species' },
        { title: '品种', dataIndex: 'breed', key: 'breed' },
        { title: '年龄', dataIndex: 'age', key: 'age', render: (a: number) =>`${a}岁` },
        {
            title: '状态', dataIndex: 'status', key: 'status',
            render: (s: string) => {
                const item = statusMap[s] || { color: 'default', text: s };
                return <Tag color={item.color}>{item.text}</Tag>;
            },
        },
        {
            title: '操作', key: 'action',
            render: (_: any, record: any) => (
                <div style={{ display: 'flex', gap: 8 }}>
                    <Button size="small" icon={<EditOutlined />} onClick={() =>openModal(record)}>
                        编辑
                    </Button>
                    <Button size="small" icon={<UploadOutlined />} onClick={() =>handleUploadImage(record.id)}>
                        传图
                    </Button>
                    <Button size="small" danger icon={<DeleteOutlined />}onClick={() => handleDelete(record.id, record.name)}>
                        删除
                    </Button>
                </div>
            ),
        },
    ];

    return (
        <>
            <div style={{ marginBottom: 16 }}>
                <Button type="primary" icon={<PlusOutlined />} onClick={() =>openModal()}>
                    添加宠物
                </Button>
            </div>

            <Table
                dataSource={pets}
                columns={columns}
                rowKey="id"
                loading={loading}
                pagination={{
                    current: page, pageSize: 10, total,
                    onChange: (p) => setPage(p),
                }}
            />

            {/* 新增/编辑弹窗 */}
            <Modal
                title={editingPet ? '编辑宠物' : '添加宠物'}
                open={modalOpen}
                onCancel={() => setModalOpen(false)}
                footer={null}
                width={520}
            >
                <Form form={form} layout="vertical" onFinish={onFinish}>
                    <Form.Item name="name" label="名字" rules={[{ required: true}]}>
                        <Input placeholder="宠物名字" />
                    </Form.Item>
                    <Form.Item name="species" label="物种" initialValue="cat">
                        <Select>
                            <Select.Option value="cat">{t('pet.cat')}</Select.Option>
                            <Select.Option value="dog">{t('pet.dog')}</Select.Option>
                            <Select.Option value="other">{t('pet.other')}</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item name="breed" label="品种">
                        <Input placeholder="如：英短、金毛" />
                    </Form.Item>
                    <Form.Item name="age" label="年龄" rules={[{ required: true }]}>
                        <InputNumber min={0} max={30} style={{ width: '100%' }} />
                    </Form.Item>
                    <Form.Item name="gender" label="性别" rules={[{ required: true}]} initialValue="male">
                        <Select>
                            <Select.Option value="male">{t('pet.male')}</Select.Option>
                            <Select.Option value="female">{t('pet.female')}</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item name="description" label="描述" rules={[{ required:true }]}>
                        <Input.TextArea rows={3} placeholder="性格、健康状况、特殊需求等" />
                    </Form.Item>
                    {editingPet && (
                        <Form.Item name="status" label="状态">
                            <Select>
                                <Select.Option value="adoptable">{t('pet.available')}</Select.Option>
                                <Select.Option value="reserved">{t('pet.reserved')}</Select.Option>
                                <Select.Option value="adopted">{t('pet.adopted')}</Select.Option>
                            </Select>
                        </Form.Item>
                    )}
                    <Form.Item>
                        <Button type="primary" htmlType="submit" loading={submitting} block>
                            {editingPet ? '保存修改' : '添加'}
                        </Button>
                    </Form.Item>
                </Form>
            </Modal>
        </>
    );
};

// ========== 领养审核 ==========
const AdoptionReviewTab = () => {
    const [adoptions, setAdoptions] = useState<any[]>([]);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);
    const [loading, setLoading] = useState(false);
    const [statusFilter, setStatusFilter] = useState('');

    const fetchAdoptions = async () => {
        setLoading(true);
        try {
            const params: Record<string, any> = { page, page_size: 10 };
            if (statusFilter) params.status = statusFilter;
            const res: any = await getOrgAdoptions(params);
            setAdoptions(res.data.list || []);
            setTotal(res.data.total || 0);
        } catch {
            //
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => { fetchAdoptions(); }, [page, statusFilter]);

    // 通过申请
    const handleApprove = (id: number) => {
        Modal.confirm({
            title: '通过领养申请',
            content: '确定通过此领养申请吗？宠物状态将变为"已预定"。',
            onOk: async () => {
                await reviewAdoption(id, { status: 'approved' });
                message.success('已通过');
                fetchAdoptions();
            },
        });
    };

    // 拒绝申请
    const handleReject = (id: number) => {
        let reason = '';
        Modal.confirm({
            title: '拒绝原因',
            content: (
                <Input.TextArea
                    rows={3}
                    placeholder="请告知申请人拒绝原因"
                    onChange={(e) => { reason = e.target.value; }}
                />
            ),
            okType: 'danger',
            onOk: async () => {
                await reviewAdoption(id, { status: 'rejected', reject_reason: reason});
                message.success('已拒绝');
                fetchAdoptions();
            },
        });
    };

    // 确认完成（已经线下交接）
    const handleComplete = (id: number) => {
        Modal.confirm({
            title: '确认领养完成',
            content: '确认宠物已完成线下交接吗？',
            onOk: async () => {
                await completeAdoption(id);
                message.success('领养已完成');
                fetchAdoptions();
            },
        });
    };

    const statusMap: Record<string, { color: string; text: string }> = {
        pending: { color: 'orange', text: '待审核' },
        approved: { color: 'green', text: '已通过' },
        rejected: { color: 'red', text: '已拒绝' },
        completed: { color: 'blue', text: '已完成' },
    };

    const columns = [
        { title: '申请人', dataIndex: ['user', 'nickname'], key: 'user', render: (n:string) => n || '-' },
        { title: '宠物', dataIndex: ['pet', 'name'], key: 'pet', render: (n: string)=> n || '-' },
        {
            title: '状态', dataIndex: 'status', key: 'status',
            render: (s: string) => {
                const item = statusMap[s] || { color: 'default', text: s };
                return <Tag color={item.color}>{item.text}</Tag>;
            },
        },
        {
            title: '申请时间', dataIndex: 'created_at', key: 'time',
            render: (t: string) => new Date(t).toLocaleDateString(),
        },
        {
            title: '操作', key: 'action',
            render: (_: any, record: any) => {
                if (record.status === 'pending') {
                    return (
                        <div style={{ display: 'flex', gap: 8 }}>
                            <Button size="small" type="primary" onClick={() =>handleApprove(record.id)}>通过</Button>
                            <Button size="small" danger onClick={() =>handleReject(record.id)}>拒绝</Button>
                        </div>
                    );
                }
                if (record.status === 'approved') {
                    return <Button size="small" onClick={() =>handleComplete(record.id)}>确认完成</Button>;
                }
                return '-';
            },
        },
    ];

    return (
        <>
            <div style={{ marginBottom: 16 }}>
                <span>状态筛选：</span>
                <Select
                    value={statusFilter}
                    onChange={(v) => { setStatusFilter(v); setPage(1); }}
                    style={{ width: 120 }}
                >
                    <Select.Option value="">全部</Select.Option>
                    <Select.Option value="pending">待审核</Select.Option>
                    <Select.Option value="approved">已通过</Select.Option>
                    <Select.Option value="rejected">已拒绝</Select.Option>
                    <Select.Option value="completed">已完成</Select.Option>
                </Select>
            </div>

            <Table
                dataSource={adoptions}
                columns={columns}
                rowKey="id"
                loading={loading}
                pagination={{
                    current: page, pageSize: 10, total,
                    onChange: (p) => setPage(p),
                }}
            />
        </>
    );
};

// ========== 主组件 ==========
const OrgDashboard = () => {
    const { t } = useTranslation();
    const [org, setOrg] = useState<any>(null);
    const [orgStatus, setOrgStatus] = useState<'none' | 'pending' | 'approved' |
'loading'>('loading');

    const fetchOrg = async () => {
        try {
            const res: any = await getMyOrg();
            setOrg(res.data);
            setOrgStatus(res.data.status === 'approved' ? 'approved' : 'pending');
        } catch {
            setOrgStatus('none');
        }
    };

    useEffect(() => { fetchOrg(); }, []);

    // 还没注册机构
    if (orgStatus === 'loading') return null;

    if (orgStatus === 'none') {
        return <OrgRegisterForm onSuccess={fetchOrg} />;
    }

    if (orgStatus === 'pending') {
        return (
            <Card style={{ textAlign: 'center', padding: 40 }}>
                <h3>⏳ {t('profile.orgPending')}</h3>
                <p>机构「{org?.name}」正在等待管理员审核</p>
                <Button onClick={fetchOrg}>刷新状态</Button>
            </Card>
        );
    }

    // 已通过 → 显示工作台
    return (
        <Tabs items={[
            {
                key: 'info',
                label: '机构信息',
                children: <OrgInfoTab org={org} onUpdate={fetchOrg} />,
            },
            {
                key: 'pets',
                label: '宠物管理',
                children: <PetManageTab orgId={org.id} />,
            },
            {
                key: 'adoptions',
                label: '领养审核',
                children: <AdoptionReviewTab />,
            },
        ]} />
    );
};

export default OrgDashboard;
