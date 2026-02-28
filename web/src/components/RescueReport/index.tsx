import { Modal, Form, Input, Select, Button, message } from 'antd';
import { useTranslation } from 'react-i18next';
import { createRescue } from '../../api/rescue';

interface RescueReportProps {
    open: boolean;
    onClose: () => void;
    onSuccess: () => void;
}

const RescueReport = ({ open, onClose, onSuccess }: RescueReportProps) => {
    const { t } = useTranslation();
    const [form] = Form.useForm();

    const onFinish = async (values: any) => {
        try {
            await createRescue({
                title: values.title,
                description: values.description,
                urgency: values.urgency,
                longitude: values.longitude ? Number(values.longitude) : 0,
                latitude: values.latitude ? Number(values.latitude) : 0,
                address: values.address || '',
            });
            message.success('上报成功');
            form.resetFields();
            onSuccess();
        } catch {
            // 拦截器已处理
        }
    };

    // 获取当前位置自动填入
    const fillCurrentLocation = () => {
        if (!navigator.geolocation) {
            message.warning('浏览器不支持定位');
            return;
        }
        navigator.geolocation.getCurrentPosition(
            (pos) => {
                form.setFieldsValue({
                    longitude: pos.coords.longitude.toFixed(6),
                    latitude: pos.coords.latitude.toFixed(6),
                });
                message.success('已获取当前位置');
            },
            () => {
                message.error('获取位置失败，请手动输入');
            }
        );
    };

    return (
        <Modal
            title={t('rescue.report')}
            open={open}
            onCancel={onClose}
            footer={null}
            width={520}
        >
            <Form form={form} layout="vertical" onFinish={onFinish}>
                <Form.Item
                    name="title"
                    label="标题"
                    rules={[{ required: true, message: '请输入标题' }]}
                >
                    <Input placeholder="如：xx路发现受伤流浪猫" />
                </Form.Item>

                <Form.Item
                    name="description"
                    label="详细描述"
                    rules={[{ required: true, message: '请描述情况' }]}
                >
                    <Input.TextArea rows={3} placeholder="描述动物状态、具体位置等"
/>
                </Form.Item>

                <Form.Item
                    name="urgency"
                    label={t('rescue.urgency')}
                    rules={[{ required: true }]}
                    initialValue="medium"
                >
                    <Select>
                        <Select.Option
value="critical">{t('rescue.critical')}</Select.Option>
                        <Select.Option
value="high">{t('rescue.high')}</Select.Option>
                        <Select.Option
value="medium">{t('rescue.medium')}</Select.Option>
                        <Select.Option value="low">{t('rescue.low')}</Select.Option>
                    </Select>
                </Form.Item>

                <Form.Item label="位置坐标">
                    <div style={{ display: 'flex', gap: 8 }}>
                        <Form.Item name="longitude" noStyle>
                            <Input placeholder="经度" style={{ flex: 1 }} />
                        </Form.Item>
                        <Form.Item name="latitude" noStyle>
                            <Input placeholder="纬度" style={{ flex: 1 }} />
                        </Form.Item>
                        <Button onClick={fillCurrentLocation}>定位</Button>
                    </div>
                </Form.Item>

                <Form.Item name="address" label="地址描述">
                    <Input placeholder="如：朝阳区xx路xx号附近" />
                </Form.Item>

                <Form.Item>
                    <Button type="primary" htmlType="submit" block>
                        提交上报
                    </Button>
                </Form.Item>
            </Form>
        </Modal>
    );
};

export default RescueReport;

