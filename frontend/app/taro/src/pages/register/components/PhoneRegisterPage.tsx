import {useState, useEffect} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Input, Button, Text} from '@tarojs/components';export default function PhoneRegisterPage() {
  const {t} = useTranslation();

  const [phone, setPhone] = useState('');
  const [verificationCode, setVerificationCode] = useState('');
  const [codeSent, setCodeSent] = useState(false);
  const [codeCountdown, setCodeCountdown] = useState(0);

  useEffect(() => {
    let codeTimer: number | null = null;

    if (codeSent && codeCountdown > 0) {
      codeTimer = window.setInterval(() => {
        setCodeCountdown((prev) => {
          if (prev <= 1) {
            setCodeSent(false);
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    }

    return () => {
      if (codeTimer !== null) {
        clearInterval(codeTimer);
      }
    };
  }, [codeSent, codeCountdown]);

  // 手机号简单校验
  const isValidPhone = () => {
    if (!phone) return false;
    const phoneRegex = /^1[3-9]\d{9}$/;
    return phoneRegex.test(phone);
  };

  /**
   * 发送验证码
   */
  const handleButtonSendVerifyCode = () => {
    if (!phone) {
      return;
    }

    setCodeSent(true);
    setCodeCountdown(60);
  };

  /**
   * 登录或者注册
   */
  const handleButtonRegisterOrLogin = () => {
    if (!phone || !verificationCode) {
      return;
    }
    console.log('手机号注册/登录:', {phone, verificationCode});
  };

  return (
    <View className='register-form'>
      {/* Phone Number Group */}
      <View className='form-group'>
        <Text className='form-label'>{t('authentication.register.phone')}</Text>
        <Input
          type='text'
          value={phone}
          onInput={(e) => setPhone(e.detail.value)}
          placeholder={t('authentication.register.input_phone')}
          className='input-field'
          maxlength={11}
        />
      </View>

      {/* Verification Code Group */}
      <View className='form-group'>
        <Text className='form-label'>{t('authentication.register.code')}</Text>
        <View className='code-input-row'>
          <Input
            type='text'
            value={verificationCode}
            onInput={(e) => setVerificationCode(e.detail.value)}
            placeholder={t('authentication.register.input_code')}
            className='input-field'
            maxlength={6}
          />
          <Button
            className='send-code-btn'
            disabled={codeSent || !isValidPhone()}
            onClick={handleButtonSendVerifyCode}
          >
            {codeSent
              ? t('authentication.register.resend_after_countdown', {count: codeCountdown})
              : t('authentication.register.send_code')}
          </Button>
        </View>
      </View>

      {/* Register/Login Button */}
      <Button
        className='register-button'
        disabled={
          !isValidPhone() ||
          !verificationCode ||
          verificationCode.length !== 6
        }
        onClick={handleButtonRegisterOrLogin}
      >
        {t('authentication.register.register_or_login')}
      </Button>
    </View>
  );
}
