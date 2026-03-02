import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, Row, Col, Tag, Button, Image, Descriptions, Spin, Empty, message, Modal, Form, Input } from 'antd';
import { HeartOutlined, HeartFilled, GiftOutlined, MessageOutlined, BookOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { getPetDetail } from '../../api/pet';
import { createAdoption } from '../../api/adoption';
import { toggleFavorite, getFavoriteStatus } from '../../api/favorite';
import useAuthStore from '../../store/useAuthStore';

interface PetImage {
    id: number;
    image_url: string;
    sort_order: number;
}

interface PetDetail {
    id: number;
    name: string;
    species: string;
    breed: string;
    age: number;
    gender: string;
    status: string;
    description: string;
    cover_image: string;
    org_id: number;
    images: PetImage[];
    organization?: {
        id: number;
        name: string;
    };
}

const statusColorMap: Record<string, string> = {
    adoptable: 'green',
    reserved: 'orange',
    adopted: 'blue',
};

const PetDetailPage = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const { t } = useTranslation();
    const { isLoggedIn } = useAuthStore();

    const [pet, setPet] = useState<PetDetail | null>(null);
    const [loading, setLoading] = useState(true);
    const [adoptLoading, setAdoptLoading] = useState(false);
    const [adoptModalOpen, setAdoptModalOpen] = useState(false);
    const [adoptForm] = Form.useForm();
    const [favorited, setFavorited] = useState(false);

    useEffect(() => {
        if (!id) return;
        setLoading(true);
        getPetDetail(Number(id))
            .then((res: any) => {
                setPet(res.data);
            })
            .catch(() => {
                setPet(null);
            })
            .finally(() => {
                setLoading(false);
            });

        // 获取收藏状态
        if (isLoggedIn) {
            getFavoriteStatus(Number(id))
                .then((res: any) => {
                    setFavorited(res.data?.favorited || false);
                })
                .catch(() => {});
        }
    }, [id, isLoggedIn]);

    // 收藏切换
    const handleToggleFavorite = async () => {
        if (!isLoggedIn) {
            message.warning('请先登录');
            return;
        }
        try {
            const res: any = await toggleFavorite(Number(id));
            setFavorited(res.data.favorited);
            message.success(res.data.favorited ? '已收藏' : '已取消收藏');
        } catch {
            //
        }
    };

    // 申请领养
    const handleAdopt = () => {
        if (!isLoggedIn) {
            message.warning('请先登录');
            navigate('/login');
            return;
        }
        setAdoptModalOpen(true);
    };

    const onAdoptSubmit = async (values: { reason: string; living_condition: string; experience?: string }) => {
        setAdoptLoading(true);
        try {
            await createAdoption({
                pet_id: Number(id),
                reason: values.reason,
                living_condition: values.living_condition,
                experience: values.experience,
            });
            message.success('申请已提交，请等待机构审核');
            setAdoptModalOpen(false);
            adoptForm.resetFields();
        } catch {
            // 拦截器已处理
        } finally {
            setAdoptLoading(false);
        }
    };

    if (loading) {
        return (
            <div style={{ textAlign: 'center', padding: 100 }}>
                <Spin size="large" />
            </div>
        );
    }

    if (!pet) {
        return <Empty description={t('pet.noPetFound')} />;
    }

    return (
        <div>
            <Card>
                <Row gutter={32}>
                    {/* 左侧：主图 */}
                    <Col xs={24} md={10}>
                        <Image
                            src={pet.cover_image || '/placeholder.png'}
                            alt={pet.name}
                            style={{ width: '100%', borderRadius: 8 }}
                            fallback="/placeholder.png"
                        />
                    </Col>

                    {/* 右侧：信息 */}
                    <Col xs={24} md={14}>
                        <div style={{ display: 'flex', alignItems: 'center',marginBottom: 16 }}>
                            <h1 style={{ margin: 0, marginRight: 1}}>{pet.name}</h1>
                            <Tag color={statusColorMap[pet.status] || 'default'}>
                                {t(`pet.${pet.status === 'adoptable' ? 'available' :pet.status}`)}
                            </Tag>
                        </div>

                        <Descriptions column={2} style={{ marginBottom: 24 }}>
                            <Descriptions.Item label={t('pet.species')}>
                                {t(`pet.${pet.species}`) || pet.species}
                            </Descriptions.Item>
                            <Descriptions.Item label={t('pet.breed')}>
                                {pet.breed || '-'}
                            </Descriptions.Item>
                            <Descriptions.Item label={t('pet.age')}>
                                {pet.age} {t('pet.ageUnit')}
                            </Descriptions.Item>
                            <Descriptions.Item label={t('pet.gender')}>
                                {t(`pet.${pet.gender}`)}
                            </Descriptions.Item>
                            {pet.organization && (
                                <Descriptions.Item label={t('pet.org')}>
                                    {pet.organization.name}
                                </Descriptions.Item>
                            )}
                        </Descriptions>

                        <h3>{t('pet.description')}</h3>
                        <p style={{ color: '#666', lineHeight: 1.8 }}>
                            {pet.description || '暂无描述'}
                        </p>

                        {/* 操作按钮 */}
                        {pet.organization && (
                                <Button size="large" icon={<MessageOutlined />}onClick={() => navigate(`/chat?userId=${pet.org_id}&name=${encodeURIComponent(pet.organization.name)}`)}
                                    style={{ marginRight: 16 }}
                                >
                                    联系机构
                                </Button>
                            )}
                        <div style={{ marginTop: 32 }}>
                            {pet.status === 'adoptable' && (
                                <Button
                                    type="primary"
                                    size="large"
                                    icon={<HeartOutlined />}
                                    loading={adoptLoading}
                                    onClick={handleAdopt}
                                    style={{ marginRight: 16 }}
                                >
                                    {t('pet.adopt')}
                                </Button>
                            )}
                            <Button
                                size="large"
                                icon={<GiftOutlined />}
                                onClick={() => navigate(`/donation?target_type=pet&target_id=${pet.id}`)}
                            >
                                {t('pet.donateForPet')}
                            </Button>
                            <Button
                                size="large"
                                icon={favorited ? <HeartFilled style={{ color: '#ff4d4f' }} /> : <HeartOutlined />}
                                onClick={handleToggleFavorite}
                                style={{ marginLeft: 16 }}
                            >
                                {favorited ? '已收藏' : '收藏'}
                            </Button>
                            {pet.status === 'adopted' && (
                                <Button
                                    size="large"
                                    icon={<BookOutlined />}
                                    onClick={() => navigate(`/pets/${pet.id}/diary`)}
                                    style={{ marginLeft: 16 }}
                                >
                                    成长日记
                                </Button>
                            )}
                            
                        </div>
                    </Col>
                </Row>
            </Card>

            {/* 更多图片 */}
            {pet.images && pet.images.length > 0 && (
                <Card title={t('pet.images')} style={{ marginTop: 24 }}>
                    <Image.PreviewGroup>
                        <Row gutter={[16, 16]}>
                            {pet.images.map((img) => (
                                <Col xs={12} sm={8} md={6} key={img.id}>
                                    <Image
                                        src={img.image_url}
                                        alt={pet.name}
                                        style={{
                                            width: '100%',
                                            height: 160,
                                            objectFit: 'cover',
                                            borderRadius: 4,
                                        }}
                                    />
                                </Col>
                            ))}
                        </Row>
                    </Image.PreviewGroup>
                </Card>
            )}

            {/* 领养申请弹窗 */}
            <Modal
                title={`申请领养「${pet.name}」`}
                open={adoptModalOpen}
                onCancel={() => setAdoptModalOpen(false)}
                footer={null}
                width={480}
            >
                <Form form={adoptForm} layout="vertical" onFinish={onAdoptSubmit}>
                    <Form.Item
                        name="reason"
                        label="领养原因"
                        rules={[{ required: true, message: '请填写领养原因' }]}
                    >
                        <Input.TextArea rows={3} placeholder="为什么想领养这只宠物？" />
                    </Form.Item>
                    <Form.Item
                        name="living_condition"
                        label="居住条件"
                        rules={[{ required: true, message: '请描述居住条件' }, { max: 200, message: '最多200字' }]}
                    >
                        <Input.TextArea rows={2} placeholder="如：有独立住房，有阳台，小区允许养宠物等" />
                    </Form.Item>
                    <Form.Item name="experience" label="养宠经验（选填）">
                        <Input.TextArea rows={2} placeholder="是否有过养宠物的经历？" />
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit" loading={adoptLoading} block>
                            提交申请
                        </Button>
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
};

export default PetDetailPage;
