import {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Input, Button, Text} from '@tarojs/components';export default function EmailRegisterPage() {
  const {t} = useTranslation();

  const [email, setEmail] = useState('');
  const [visibleEnter, setVisibleEnter] = useState(false);

  // 简单的邮箱格式验证
  const isValidEmail = () => {
    if (!email) return false;
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  const handleButtonNext = () => {
    if (!isValidEmail()) {
      return;
    }
    setVisibleEnter(true);
  };

  return (
    <View>
      {!visibleEnter ? (
        <View className='register-form'>
          {/* Email Input Group */}
          <View className='form-group'>
            <Text className='form-label'>{t('authentication.register.email')}</Text>
            <Input
              type='text'
              value={email}
              onInput={e => setEmail(e.detail.value)}
              placeholder={t('authentication.register.input_email')}
              className={`input-field${email && !isValidEmail() ? ' error' : ''}`}
            />
            {email && !isValidEmail() && (
              <Text className='error-hint'>{t('authentication.register.invalid_email')}</Text>
            )}
          </View>
          {/* Next Button */}
          <Button
            className='register-button'
            disabled={!isValidEmail()}
            onClick={handleButtonNext}
          >
            {t('authentication.register.next_step')}
          </Button>
        </View>
      ) : (
        <EmailRegisterEnterCodePage email={email} />
      )}
    </View>
  );
}

function EmailRegisterEnterCodePage({email}: { email: string }) {
  const {t} = useTranslation();
  const [code, setCode] = useState('');
  const [codeSent, setCodeSent] = useState(true);
  const [codeCountdown, setCodeCountdown] = useState(60);

  // 启动倒计时
  const startCountdown = () => {
    setCodeSent(true);
    setCodeCountdown(60);
    const timer = setInterval(() => {
      setCodeCountdown(prev => {
        if (prev <= 1) {
          clearInterval(timer);
          setCodeSent(false);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
  };

  const handleResend = () => {
    if (!codeSent) {
      startCountdown();
    }
  };

  const handleRegister = () => {
    if (!code || code.length !== 6) return;
    // TODO: 注册逻辑
  };

  return (
    <View className='register-form'>
      <View className='email-sent-info'>
        <Text className='hint-title'>{t('authentication.register.code_sent_title')}</Text>
        <Text className='email-display'>{email}</Text>
        <Text className='hint-subtitle'>{t('authentication.register.code_sent_subtitle')}</Text>
      </View>
      <View className='form-group'>
        <Text className='form-label'>{t('authentication.register.code')}</Text>
        <Input
          type='text'
          value={code}
          onInput={e => setCode(e.detail.value)}
          placeholder={t('authentication.register.input_code')}
          className='input-field'
          maxlength={6}
        />
      </View>
      <Button
        className='register-button'
        disabled={!code || code.length !== 6}
        onClick={handleRegister}
      >
        {t('authentication.register.register')}
      </Button>
      <View className='resend-section'>
        <Text>
          {codeSent
            ? t('authentication.register.resend_after_countdown', {count: codeCountdown})
            : (
              <Text className='text-btn' onClick={handleResend}>
                {t('authentication.register.resend_code')}
              </Text>
            )}
        </Text>
      </View>
    </View>
  );
}
