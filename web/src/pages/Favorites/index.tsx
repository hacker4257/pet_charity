import { useEffect, useState } from 'react';
import { Card, Row, Col, Empty, Spin, Tag, Button } from 'antd';
import { HeartFilled } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { getMyFavorites } from '../../api/favorite';
import { toggleFavorite } from '../../api/favorite';

const Favorites = () => {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const [list, setList] = useState<any[]>([]);
    const [loading, setLoading] = useState(false);

    const fetchList = async () => {
        setLoading(true);
        try {
            const res: any = await getMyFavorites();
            setList(res.data || []);
        } catch {
            //
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchList();
    }, []);

    const handleUnfavorite = async (petId: number) => {
        try {
            await toggleFavorite(petId);
            setList(prev => prev.filter(item => item.id !== petId));
        } catch {
            //
        }
    };

    if (loading) {
        return <div style={{ textAlign: 'center', padding: 60 }}><Spin size="large" /></div>;
    }

    return (
        <div style={{ maxWidth: 1200, margin: '0 auto' }}>
            <h2><HeartFilled style={{ color: '#ff4d4f' }} /> {t('nav.favorites') || '我的收藏'}</h2>

            {list.length === 0 ? (
                <Empty description="暂无收藏的宠物" />
            ) : (
                <Row gutter={[16, 16]}>
                    {list.map((pet: any) => (
                        <Col xs={24} sm={12} md={8} lg={6} key={pet.id}>
                            <Card
                                hoverable
                                cover={
                                    pet.cover_image ? (
                                        <img
                                            alt={pet.name}
                                            src={pet.cover_image}
                                            style={{ height: 200, objectFit: 'cover' }}
                                        />
                                    ) : (
                                        <div style={{
                                            height: 200,
                                            background: '#f0f0f0',
                                            display: 'flex',
                                            alignItems: 'center',
                                            justifyContent: 'center',
                                            color: '#999',
                                        }}>
                                            暂无图片
                                        </div>
                                    )
                                }
                                onClick={() => navigate(`/pets/${pet.id}`)}
                                actions={[
                                    <Button
                                        key="unfav"
                                        type="text"
                                        danger
                                        icon={<HeartFilled />}
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            handleUnfavorite(pet.id);
                                        }}
                                    >
                                        取消收藏
                                    </Button>,
                                ]}
                            >
                                <Card.Meta
                                    title={pet.name}
                                    description={
                                        <div>
                                            <Tag>{pet.species}</Tag>
                                            {pet.breed && <Tag>{pet.breed}</Tag>}
                                            <Tag color={pet.gender === 'male' ? 'blue' : 'pink'}>
                                                {pet.gender === 'male' ? '公' : '母'}
                                            </Tag>
                                        </div>
                                    }
                                />
                            </Card>
                        </Col>
                    ))}
                </Row>
            )}
        </div>
    );
};

export default Favorites;
