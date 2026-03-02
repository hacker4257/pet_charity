import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
    Card, Row, Col, Button, Image, Spin, Empty, message,
    Modal, Form, Input, Pagination, Avatar, Upload,
} from 'antd';
import {
    HeartOutlined, HeartFilled, EditOutlined, DeleteOutlined,
    PlusOutlined, CameraOutlined, ArrowLeftOutlined,
} from '@ant-design/icons';
import {
    getDiariesByPet, createDiary, updateDiary, deleteDiary,
    uploadDiaryImage, deleteDiaryImage, toggleDiaryLike,
} from '../../api/diary';
import { getPetDetail } from '../../api/pet';
import useAuthStore from '../../store/useAuthStore';

interface DiaryImage {
    id: number;
    image_url: string;
    sort_order: number;
}

interface DiaryItem {
    id: number;
    user_id: number;
    pet_id: number;
    content: string;
    created_at: string;
    updated_at: string;
    user?: { id: number; nickname: string; avatar: string };
    images?: DiaryImage[];
    liked?: boolean;
    like_count?: number;
}

interface PetInfo {
    id: number;
    name: string;
    cover_image: string;
    species: string;
}

const PetDiary = () => {
    const { petId } = useParams<{ petId: string }>();
    const navigate = useNavigate();
    const { isLoggedIn, user } = useAuthStore();

    const [pet, setPet] = useState<PetInfo | null>(null);
    const [diaries, setDiaries] = useState<DiaryItem[]>([]);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);
    const [loading, setLoading] = useState(true);

    // 创建/编辑弹窗
    const [modalOpen, setModalOpen] = useState(false);
    const [editingDiary, setEditingDiary] = useState<DiaryItem | null>(null);
    const [form] = Form.useForm();
    const [submitting, setSubmitting] = useState(false);

    // 图片上传弹窗
    const [imageModalOpen, setImageModalOpen] = useState(false);
    const [imageDiaryId, setImageDiaryId] = useState<number | null>(null);
    const [uploading, setUploading] = useState(false);

    const pageSize = 10;

    // 加载宠物信息
    useEffect(() => {
        if (!petId) return;
        getPetDetail(Number(petId))
            .then((res: any) => setPet(res.data))
            .catch(() => setPet(null));
    }, [petId]);

    // 加载日记列表
    const fetchDiaries = () => {
        if (!petId) return;
        setLoading(true);
        getDiariesByPet(Number(petId), { page, page_size: pageSize })
            .then((res: any) => {
                setDiaries(res.data?.items || []);
                setTotal(res.data?.total || 0);
            })
            .catch(() => {
                setDiaries([]);
            })
            .finally(() => setLoading(false));
    };

    useEffect(() => {
        fetchDiaries();
    }, [petId, page]);

    // 创建日记
    const handleCreate = () => {
        if (!isLoggedIn) {
            message.warning('请先登录');
            navigate('/login');
            return;
        }
        setEditingDiary(null);
        form.resetFields();
        setModalOpen(true);
    };

    // 编辑日记
    const handleEdit = (diary: DiaryItem) => {
        setEditingDiary(diary);
        form.setFieldsValue({ content: diary.content });
        setModalOpen(true);
    };

    // 提交创建/编辑
    const onSubmit = async (values: { content: string }) => {
        setSubmitting(true);
        try {
            if (editingDiary) {
                await updateDiary(editingDiary.id, { content: values.content });
                message.success('日记已更新');
            } else {
                await createDiary({ pet_id: Number(petId), content: values.content });
                message.success('日记已发布');
            }
            setModalOpen(false);
            form.resetFields();
            fetchDiaries();
        } catch {
            // 拦截器已处理
        } finally {
            setSubmitting(false);
        }
    };

    // 删除日记
    const handleDelete = (diaryId: number) => {
        Modal.confirm({
            title: '确认删除',
            content: '删除后不可恢复，确定要删除这篇日记吗？',
            okType: 'danger',
            onOk: async () => {
                try {
                    await deleteDiary(diaryId);
                    message.success('已删除');
                    fetchDiaries();
                } catch {
                    // 拦截器已处理
                }
            },
        });
    };

    // 点赞
    const handleLike = async (diaryId: number) => {
        if (!isLoggedIn) {
            message.warning('请先登录');
            return;
        }
        try {
            const res: any = await toggleDiaryLike(diaryId);
            const liked = res.data?.liked;
            setDiaries(prev =>
                prev.map(d =>
                    d.id === diaryId
                        ? {
                              ...d,
                              liked,
                              like_count: (d.like_count || 0) + (liked ? 1 : -1),
                          }
                        : d
                )
            );
        } catch {
            // 拦截器已处理
        }
    };

    // 上传图片
    const openImageUpload = (diaryId: number) => {
        setImageDiaryId(diaryId);
        setImageModalOpen(true);
    };

    const handleImageUpload = async (file: File) => {
        if (!imageDiaryId) return false;
        setUploading(true);
        try {
            await uploadDiaryImage(imageDiaryId, file);
            message.success('图片已上传');
            setImageModalOpen(false);
            fetchDiaries();
        } catch {
            // 拦截器已处理
        } finally {
            setUploading(false);
        }
        return false; // 阻止 antd 自动上传
    };

    // 删除图片
    const handleDeleteImage = async (diaryId: number, imageId: number) => {
        try {
            await deleteDiaryImage(diaryId, imageId);
            message.success('图片已删除');
            fetchDiaries();
        } catch {
            // 拦截器已处理
        }
    };

    // 格式化时间
    const formatDate = (dateStr: string) => {
        const d = new Date(dateStr);
        return d.toLocaleDateString('zh-CN', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
        });
    };

    const isOwner = (diary: DiaryItem) => user && user.id === diary.user_id;

    return (
        <div style={{ maxWidth: 800, margin: '0 auto', padding: '24px 16px' }}>
            {/* 头部 */}
            <div style={{ marginBottom: 24 }}>
                <Button
                    type="text"
                    icon={<ArrowLeftOutlined />}
                    onClick={() => navigate(`/pets/${petId}`)}
                    style={{ marginBottom: 16 }}
                >
                    返回宠物详情
                </Button>

                {pet && (
                    <Card>
                        <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
                            <Avatar
                                src={pet.cover_image || '/placeholder.png'}
                                size={64}
                                shape="square"
                                style={{ borderRadius: 8 }}
                            />
                            <div style={{ flex: 1 }}>
                                <h2 style={{ margin: 0 }}>{pet.name} 的成长日记</h2>
                                <span style={{ color: '#999' }}>
                                    共 {total} 篇日记
                                </span>
                            </div>
                            <Button
                                type="primary"
                                icon={<PlusOutlined />}
                                onClick={handleCreate}
                            >
                                写日记
                            </Button>
                        </div>
                    </Card>
                )}
            </div>

            {/* 日记列表 */}
            {loading ? (
                <div style={{ textAlign: 'center', padding: 100 }}>
                    <Spin size="large" />
                </div>
            ) : diaries.length === 0 ? (
                <Empty description="还没有日记，快来写第一篇吧">
                    <Button type="primary" onClick={handleCreate}>
                        写日记
                    </Button>
                </Empty>
            ) : (
                <>
                    {diaries.map(diary => (
                        <Card
                            key={diary.id}
                            style={{ marginBottom: 16 }}
                            actions={[
                                <span key="like" onClick={() => handleLike(diary.id)}>
                                    {diary.liked ? (
                                        <HeartFilled style={{ color: '#ff4d4f' }} />
                                    ) : (
                                        <HeartOutlined />
                                    )}
                                    {' '}{diary.like_count || 0}
                                </span>,
                                ...(isOwner(diary)
                                    ? [
                                          <span key="image" onClick={() => openImageUpload(diary.id)}>
                                              <CameraOutlined /> 添加图片
                                          </span>,
                                          <span key="edit" onClick={() => handleEdit(diary)}>
                                              <EditOutlined /> 编辑
                                          </span>,
                                          <span key="delete" onClick={() => handleDelete(diary.id)}>
                                              <DeleteOutlined /> 删除
                                          </span>,
                                      ]
                                    : []),
                            ]}
                        >
                            {/* 作者信息 */}
                            <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 12 }}>
                                <Avatar src={diary.user?.avatar} size={36}>
                                    {diary.user?.nickname?.[0] || '?'}
                                </Avatar>
                                <div>
                                    <div style={{ fontWeight: 500 }}>
                                        {diary.user?.nickname || '匿名用户'}
                                    </div>
                                    <div style={{ fontSize: 12, color: '#999' }}>
                                        {formatDate(diary.created_at)}
                                    </div>
                                </div>
                            </div>

                            {/* 日记内容 */}
                            <p style={{ lineHeight: 1.8, whiteSpace: 'pre-wrap' }}>
                                {diary.content}
                            </p>

                            {/* 图片 */}
                            {diary.images && diary.images.length > 0 && (
                                <Image.PreviewGroup>
                                    <Row gutter={[8, 8]} style={{ marginTop: 12 }}>
                                        {diary.images.map(img => (
                                            <Col key={img.id} xs={8} sm={6}>
                                                <div style={{ position: 'relative' }}>
                                                    <Image
                                                        src={img.image_url}
                                                        style={{
                                                            width: '100%',
                                                            height: 120,
                                                            objectFit: 'cover',
                                                            borderRadius: 4,
                                                        }}
                                                    />
                                                    {isOwner(diary) && (
                                                        <Button
                                                            type="text"
                                                            size="small"
                                                            danger
                                                            icon={<DeleteOutlined />}
                                                            style={{
                                                                position: 'absolute',
                                                                top: 4,
                                                                right: 4,
                                                                background: 'rgba(255,255,255,0.8)',
                                                                borderRadius: '50%',
                                                                minWidth: 24,
                                                                padding: 0,
                                                            }}
                                                            onClick={(e) => {
                                                                e.stopPropagation();
                                                                handleDeleteImage(diary.id, img.id);
                                                            }}
                                                        />
                                                    )}
                                                </div>
                                            </Col>
                                        ))}
                                    </Row>
                                </Image.PreviewGroup>
                            )}
                        </Card>
                    ))}

                    {total > pageSize && (
                        <div style={{ textAlign: 'center', marginTop: 24 }}>
                            <Pagination
                                current={page}
                                pageSize={pageSize}
                                total={total}
                                onChange={(p) => setPage(p)}
                                showSizeChanger={false}
                            />
                        </div>
                    )}
                </>
            )}

            {/* 创建/编辑弹窗 */}
            <Modal
                title={editingDiary ? '编辑日记' : '写日记'}
                open={modalOpen}
                onCancel={() => setModalOpen(false)}
                footer={null}
                width={520}
            >
                <Form form={form} layout="vertical" onFinish={onSubmit}>
                    <Form.Item
                        name="content"
                        label="日记内容"
                        rules={[
                            { required: true, message: '请填写日记内容' },
                            { max: 2000, message: '最多2000字' },
                        ]}
                    >
                        <Input.TextArea
                            rows={6}
                            placeholder="记录宠物的成长点滴..."
                            showCount
                            maxLength={2000}
                        />
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit" loading={submitting} block>
                            {editingDiary ? '保存修改' : '发布日记'}
                        </Button>
                    </Form.Item>
                </Form>
            </Modal>

            {/* 上传图片弹窗 */}
            <Modal
                title="上传图片"
                open={imageModalOpen}
                onCancel={() => setImageModalOpen(false)}
                footer={null}
                width={400}
            >
                <Upload.Dragger
                    accept="image/*"
                    showUploadList={false}
                    beforeUpload={(file) => {
                        handleImageUpload(file);
                        return false;
                    }}
                    disabled={uploading}
                >
                    <p className="ant-upload-drag-icon">
                        <CameraOutlined style={{ fontSize: 48, color: '#999' }} />
                    </p>
                    <p>{uploading ? '上传中...' : '点击或拖拽图片到此处上传'}</p>
                </Upload.Dragger>
            </Modal>
        </div>
    );
};

export default PetDiary;
