import { useEffect, useState } from 'react';
import { Card, Row, Col, Select, Input, Pagination, Tag, Empty, Spin } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { getPets } from '../../api/pet';

interface Pet {
    id: number;
    name: string;
    species: string;
    breed: string;
    age: number;
    gender: string;
    status: string;
    cover_image: string;
    description: string;
}

const statusColorMap: Record<string, string> = {
    adoptable: 'green',
    reserved: 'orange',
    adopted: 'blue',
};

const Pets = () => {
    const navigate = useNavigate();
    const { t } = useTranslation();
    const [searchParams, setSearchParams] = useSearchParams();

    // 从 URL 读取筛选参数（刷新不丢失）
    const [species, setSpecies] = useState(searchParams.get('species') || '');
    const [gender, setGender] = useState(searchParams.get('gender') || '');
    const [status, setStatus] = useState(searchParams.get('status') || '');
    const [keyword, setKeyword] = useState(searchParams.get('keyword') || '');
    const [page, setPage] = useState(Number(searchParams.get('page')) || 1);
    const [pageSize] = useState(12);

    const [pets, setPets] = useState<Pet[]>([]);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(false);

    // 拉取数据
    const fetchPets = async () => {
        setLoading(true);
        try {
            const params: Record<string, any> = { page, pageSize };
            if (species) params.species = species;
            if (gender) params.gender = gender;
            if (status) params.status = status;

            const res: any = await getPets(params);
            setPets(res.data.list || []);
            setTotal(res.data.total || 0);
        } catch {
            // 已有拦截器处理
        } finally {
            setLoading(false);
        }
    };

    // 筛选变化时重新请求，同时同步到 URL
    useEffect(() => {
        fetchPets();

        // 把筛选条件写进 URL
        const params: Record<string, string> = {};
        if (species) params.species = species;
        if (gender) params.gender = gender;
        if (status) params.status = status;
        if (keyword) params.keyword = keyword;
        if (page > 1) params.page = String(page);
        setSearchParams(params, { replace: true });
    }, [species, gender, status, page]);

    // 筛选变化时回到第 1 页
    const handleFilterChange = (setter: (v: string) => void) => {
        return (value: string) => {
            setter(value);
            setPage(1);
        };
    };

    return (
        <div>
            {/* 筛选栏 */}
            <Card style={{ marginBottom: 24 }}>
                <Row gutter={16} align="middle">
                    <Col>
                        <span>{t('pet.species')}：</span>
                        <Select
                            value={species}
                            onChange={handleFilterChange(setSpecies)}
                            style={{ width: 120 }}
                        >
                            <Select.Option value="">{t('pet.all')}</Select.Option>
                            <Select.Option value="cat">{t('pet.cat')}</Select.Option>
                            <Select.Option value="dog">{t('pet.dog')}</Select.Option>
                            <Select.Option value="other">{t('pet.other')}</Select.Option>
                        </Select>
                    </Col>
                    <Col>
                        <span>{t('pet.gender')}：</span>
                        <Select
                            value={gender}
                            onChange={handleFilterChange(setGender)}
                            style={{ width: 120 }}
                        >
                            <Select.Option value="">{t('pet.all')}</Select.Option>
                            <Select.Option value="male">{t('pet.male')}</Select.Option>
                            <Select.Option value="female">{t('pet.female')}</Select.Option>
                        </Select>
                    </Col>
                    <Col>
                        <span>{t('pet.status')}：</span>
                        <Select
                            value={status}
                            onChange={handleFilterChange(setStatus)}
                            style={{ width: 120 }}
                        >
                            <Select.Option value="">{t('pet.all')}</Select.Option>
                            <Select.Option value="adoptable">{t('pet.available')}</Select.Option>
                            <Select.Option value="reserved">{t('pet.reserved')}</Select.Option>
                            <Select.Option value="adopted">{t('pet.adopted')}</Select.Option>
                        </Select>
                    </Col>
                    <Col flex="auto" style={{ textAlign: 'right' }}>
                        <Input
                            prefix={<SearchOutlined />}
                            placeholder={t('common.search')}
                            value={keyword}
                            onChange={(e) => setKeyword(e.target.value)}
                            onPressEnter={() => { setPage(1); fetchPets(); }}
                            style={{ width: 240 }}
                            allowClear
                        />
                    </Col>
                </Row>
            </Card>

            {/* 宠物卡片列表 */}
            <Spin spinning={loading}>
                {pets.length > 0 ? (
                    <Row gutter={[24, 24]}>
                        {pets.map((pet) => (
                            <Col xs={24} sm={12} md={8} lg={6} key={pet.id}>
                                <Card
                                    hoverable
                                    onClick={() => navigate(`/pets/${pet.id}`)}
                                    cover={
                                        <img
                                            alt={pet.name}
                                            src={pet.cover_image || '/placeholder.png'}
                                            style={{ height: 200, objectFit: 'cover'}}
                                        />
                                    }
                                >
                                    <Card.Meta
                                        title={
                                            <div style={{ display: 'flex',justifyContent: 'space-between' }}>
                                                <span>{pet.name}</span>
                                                <Tag color={statusColorMap[pet.status] || 'default'}>
                                                    {t(`pet.${pet.status === 'adoptable' ? 'available' : pet.status}`)}
                                                </Tag>
                                            </div>
                                        }
                                        description={
                                            <div>
                                                <div>{pet.breed} · {pet.age}{t('pet.ageUnit')} · {t(`pet.${pet.gender}`)}</div>
                                            </div>
                                        }
                                    />
                                </Card>
                            </Col>
                        ))}
                    </Row>
                ) : (
                    <Empty description={t('common.noData')} />
                )}
            </Spin>

            {/* 分页 */}
            {total > pageSize && (
                <div style={{ textAlign: 'center', marginTop: 32 }}>
                    <Pagination
                        current={page}
                        pageSize={pageSize}
                        total={total}
                        onChange={(p) => setPage(p)}
                        showTotal={(total) => `共 ${total} 条`}
                    />
                </div>
            )}
        </div>
    );
};

export default Pets;