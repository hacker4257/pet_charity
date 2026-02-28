import { useEffect, useState, useRef } from 'react';
import { Tabs, Card, Radio, InputNumber, Input, Button, List, Tag, Modal,
Result, message } from 'antd';
import { HeartOutlined, WechatOutlined, AlipayCircleOutlined, TrophyOutlined } from
'@ant-design/icons';
import { useSearchParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { createDonation, getDonationStatus, getPublicDonations, getMyDonations }
from '../../api/donation';
import useAuthStore from '../../store/useAuthStore';

// 预设金额（单位：元）
const PRESET_AMOUNTS = [10, 50, 100, 200, 500];

// 金额分转元
const fenToYuan = (fen: number) => (fen / 100).toFixed(2);

const Donation = () => {
    const { t } = useTranslation();
    const { isLoggedIn } = useAuthStore();
    const [searchParams] = useSearchParams();

    // 从 URL 读取（宠物详情页跳转过来时会带参数）
    const urlTargetType = searchParams.get('target_type') || 'platform';
    const urlTargetId = searchParams.get('target_id');

    // ========== Tab 1: 我要捐赠 ==========
    const [targetType, setTargetType] = useState(urlTargetType);
    const [targetId, setTargetId] = useState(urlTargetId ? Number(urlTargetId) :
undefined);
    const [amount, setAmount] = useState<number>(0);
    const [customAmount, setCustomAmount] = useState<number | null>(null);
    const [payMethod, setPayMethod] = useState('wechat');
    const [remark, setRemark] = useState('');
    const [submitting, setSubmitting] = useState(false);

    // 支付状态轮询
    const [payModalOpen, setPayModalOpen] = useState(false);
    const [payUrl, setPayUrl] = useState('');
    const [payStatus, setPayStatus] = useState<'pending' | 'paid' |'failed'>('pending');
    const pollRef = useRef<ReturnType<typeof setInterval> | null>(null);
    const donationIdRef = useRef<number>(0);

    // 选择预设金额
    const handlePresetClick = (val: number) => {
        setAmount(val);
        setCustomAmount(null);
    };

    // 输入自定义金额
    const handleCustomAmount = (val: number | null) => {
        setCustomAmount(val);
        if (val && val > 0) {
            setAmount(val);
        }
    };

    // 实际金额（元）
    const finalAmount = amount;

    // 提交捐赠
    const handleDonate = async () => {
        if (!isLoggedIn) {
            message.warning('请先登录');
            return;
        }
        if (finalAmount <= 0) {
            message.warning('请选择或输入金额');
            return;
        }
        if (targetType !== 'platform' && !targetId) {
            message.warning(t('donation.noTarget'));
            return;
        }

        setSubmitting(true);
        try {
            const res: any = await createDonation({
                target_type: targetType,
                target_id: targetType !== 'platform' ? targetId : undefined,
                amount: finalAmount * 100, // 元转分
                message: remark || undefined,
                payment_method: payMethod,
            });

            // 后端返回支付链接和订单 ID
            donationIdRef.current = res.data.donation_id;
            setPayUrl(res.data.pay_url || '');
            setPayStatus('pending');
            setPayModalOpen(true);

            // 开始轮询支付状态
            startPolling(res.data.donation_id);
        } catch {
            //
        } finally {
            setSubmitting(false);
        }
    };

    // 轮询支付状态
    const startPolling = (id: number) => {
        // 先清除旧轮询
        if (pollRef.current) clearInterval(pollRef.current);

        pollRef.current = setInterval(async () => {
            try {
                const res: any = await getDonationStatus(id);
                if (res.data.status === 'paid') {
                    setPayStatus('paid');
                    stopPolling();
                } else if (res.data.status === 'failed') {
                    setPayStatus('failed');
                    stopPolling();
                }
            } catch {
                //
            }
        }, 3000); // 每 3 秒查一次
    };

    const stopPolling = () => {
        if (pollRef.current) {
            clearInterval(pollRef.current);
            pollRef.current = null;
        }
    };

    // 组件卸载时清除轮询
    useEffect(() => {
        return () => stopPolling();
    }, []);

    // ========== Tab 2: 捐赠榜 ==========
    const [publicList, setPublicList] = useState<any[]>([]);
    const [publicTotal, setPublicTotal] = useState(0);
    const [publicPage, setPublicPage] = useState(1);

    const fetchPublicList = async () => {
        try {
            const res: any = await getPublicDonations({ page: publicPage, page_size:10 });
            setPublicList(res.data.list || []);
            setPublicTotal(res.data.total || 0);
        } catch {
            //
        }
    };

    useEffect(() => {
        fetchPublicList();
    }, [publicPage]);

    // ========== Tab 3: 我的捐赠 ==========
    const [myList, setMyList] = useState<any[]>([]);
    const [myTotal, setMyTotal] = useState(0);
    const [myPage, setMyPage] = useState(1);

    const fetchMyList = async () => {
        if (!isLoggedIn) return;
        try {
            const res: any = await getMyDonations({ page: myPage, page_size: 10 });
            setMyList(res.data.list || []);
            setMyTotal(res.data.total || 0);
        } catch {
            //
        }
    };

    useEffect(() => {
        fetchMyList();
    }, [myPage]);

    // 支付状态颜色
    const statusTag = (s: string) => {
        const map: Record<string, { color: string; label: string }> = {
            pending: { color: 'orange', label: t('donation.waitingPay') },
            paid: { color: 'green', label: t('donation.paid') },
            failed: { color: 'red', label: t('donation.failed') },
        };
        const item = map[s] || { color: 'default', label: s };
        return <Tag color={item.color}>{item.label}</Tag>;
    };

    return (
        <div style={{ maxWidth: 800, margin: '0 auto' }}>
            <Tabs
                centered
                items={[
                    {
                        key: 'donate',
                        label: <span><HeartOutlined />{t('donation.donate')}</span>,
                        children: (
                            <Card>
                                {/* 捐赠对象 */}
                                <div style={{ marginBottom: 24 }}>
                                    <h4>捐赠对象</h4>
                                    <Radio.Group
                                        value={targetType}
                                        onChange={(e) => {
                                            setTargetType(e.target.value);
                                            if (e.target.value === 'platform') setTargetId(undefined);
                                        }}
                                    >
                                        <Radio.Button value="platform">{t('donation.toPlatform')}</Radio.Button>
                                        <Radio.Button value="organization">{t('donation.toOrg')}</Radio.Button>
                                        <Radio.Button value="pet">{t('donation.toPet')}</Radio.Button>
                                    </Radio.Group>

                                    {targetType !== 'platform' && (
                                        <InputNumber
                                            style={{ marginLeft: 16, width: 200 }}
                                            placeholder={targetType === 'organization' ? '机构 ID' : '宠物 ID'}
                                            value={targetId}
                                            onChange={(v) => setTargetId(v || undefined)}
                                            min={1}
                                        />
                                    )}
                                </div>

                                {/* 金额选择 */}
                                <div style={{ marginBottom: 24 }}>
                                    <h4>{t('donation.amount')}</h4>
                                    <div style={{ display: 'flex', gap: 12,flexWrap: 'wrap' }}>
                                        {PRESET_AMOUNTS.map((val) => (
                                            <Button
                                                key={val}
                                                type={amount === val && !customAmount ? 'primary' : 'default'}
                                                size="large"
                                                onClick={() => handlePresetClick(val)}
                                                style={{ width: 100 }}
                                            >
                                                ¥{val}
                                            </Button>
                                        ))}
                                        <InputNumber
                                            size="large"
                                            placeholder={t('donation.customAmount')}
                                            value={customAmount}
                                            onChange={handleCustomAmount}
                                            min={1}
                                            max={100000}
                                            prefix="¥"
                                            style={{ width: 160 }}
                                        />
                                    </div>
                                </div>

                                {/* 支付方式 */}
                                <div style={{ marginBottom: 24 }}>
                                    <h4>{t('donation.payMethod')}</h4>
                                    <Radio.Group
                                        value={payMethod}
                                        onChange={(e) => setPayMethod(e.target.value)}
                                        size="large"
                                    >
                                        <Radio.Button value="wechat">
                                            <WechatOutlined style={{ color:'#07c160' }} /> {t('donation.wechat')}
                                        </Radio.Button>
                                        <Radio.Button value="alipay">
                                            <AlipayCircleOutlined style={{ color:'#1677ff' }} /> {t('donation.alipay')}
                                        </Radio.Button>
                                    </Radio.Group>
                                </div>

                                {/* 留言 */}
                                <div style={{ marginBottom: 24 }}>
                                    <h4>{t('donation.remark')}</h4>
                                    <Input.TextArea
                                        rows={2}
                                        value={remark}
                                        onChange={(e) => setRemark(e.target.value)}
                                        placeholder="写点什么鼓励它们吧~"
                                        maxLength={200}
                                        showCount
                                    />
                                </div>

                                {/* 提交 */}
                                <Button
                                    type="primary"
                                    size="large"
                                    block
                                    loading={submitting}
                                    onClick={handleDonate}
                                    disabled={finalAmount <= 0}
                                >
                                    {finalAmount > 0
                                        ? `${t('donation.donate')} ¥${finalAmount}`
                                        : t('donation.donate')}
                                </Button>
                            </Card>
                        ),
                    },
                    {
                        key: 'public',
                        label: <span><TrophyOutlined />{t('donation.publicList')}</span>,
                        children: (
                            <List
                                dataSource={publicList}
                                locale={{ emptyText: t('common.noData') }}
                                pagination={{
                                    current: publicPage,
                                    pageSize: 10,
                                    total: publicTotal,
                                    onChange: (p) => setPublicPage(p),
                                }}
                                renderItem={(item: any, index: number) => (
                                    <List.Item>
                                        <List.Item.Meta
                                            avatar={
                                                <div style={{
                                                    width: 36, height: 36,borderRadius: '50%',
                                                    background: index < 3 ?'#ffd700' : '#f0f0f0',
                                                    display: 'flex', alignItems:'center', justifyContent: 'center',
                                                    fontWeight: 'bold',
                                                    color: index < 3 ? '#fff' :'#999',
                                                }}>
                                                    {(publicPage - 1) * 10 + index +1}
                                                </div>
                                            }
                                            title={item.user?.nickname ||'匿名用户'}
                                            description={`捐给${
                                                item.target_type === 'platform' ?'平台'
                                                : item.target_type ==='organization' ? '救助站'
                                                : '宠物'
                                            }`}
                                        />
                                        <div style={{ fontSize: 18, color:'#f5222d', fontWeight: 'bold' }}>
                                            ¥{fenToYuan(item.amount)}
                                        </div>
                                    </List.Item>
                                )}
                            />
                        ),
                    },
                    {
                        key: 'mine',
                        label: <span>{t('donation.myDonations')}</span>,
                        children: !isLoggedIn ? (
                            <Result status="warning" title="请先登录查看" />
                        ) : (
                            <List
                                dataSource={myList}
                                locale={{ emptyText: t('common.noData') }}
                                pagination={{
                                    current: myPage,
                                    pageSize: 10,
                                    total: myTotal,
                                    onChange: (p) => setMyPage(p),
                                }}
                                renderItem={(item: any) => (
                                    <List.Item>
                                        <List.Item.Meta
                                            title={
                                                <span>
                                                    ¥{fenToYuan(item.amount)}
                                                    <span style={{ marginLeft: 12}}>{statusTag(item.status)}</span>
                                                </span>
                                            }
                                            description={
                                                <span>
                                                    {new Date(item.created_at).toLocaleString()}
                                                    {' · '}
                                                    {item.target_type === 'platform'? '平台'
                                                        : item.target_type ==='organization' ? '救助站'
                                                        : '宠物'}
                                                </span>
                                            }
                                        />
                                    </List.Item>
                                )}
                            />
                        ),
                    },
                ]}
            />

            {/* 支付弹窗 */}
            <Modal
                open={payModalOpen}
                footer={null}
                onCancel={() => {
                    stopPolling();
                    setPayModalOpen(false);
                }}
                title="完成支付"
            >
                {payStatus === 'pending' && (
                    <div style={{ textAlign: 'center', padding: 24 }}>
                        {payUrl ? (
                            <>
                                <p>请使用{payMethod === 'wechat' ? '微信' :'支付宝'}扫码支付</p>
                                {/* 微信支付返回的是二维码链接，实际项目用 qrcode库渲染 */}
                                <div style={{
                                    padding: 20, background: '#f5f5f5',
                                    borderRadius: 8, marginBottom: 16,
                                    wordBreak: 'break-all',
                                }}>
                                    <code>{payUrl}</code>
                                </div>
                                <p style={{ color: '#999'}}>支付完成后页面会自动更新...</p>
                            </>
                        ) : (
                            <p>正在创建支付订单...</p>
                        )}
                    </div>
                )}
                {payStatus === 'paid' && (
                    <Result
                        status="success"
                        title={t('donation.paySuccess')}
                        extra={
                            <Button type="primary" onClick={() => {
                                setPayModalOpen(false);
                                fetchMyList();
                            }}>
                                确定
                            </Button>
                        }
                    />
                )}
                {payStatus === 'failed' && (
                    <Result
                        status="error"
                        title="支付失败"
                        extra={
                            <Button onClick={() => setPayModalOpen(false)}>
                                关闭
                            </Button>
                        }
                    />
                )}
            </Modal>
        </div>
    );
};

export default Donation;
