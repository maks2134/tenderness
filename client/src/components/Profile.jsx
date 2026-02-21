import React, { useState, useEffect } from 'react'
import OAuthLink from './OAuthLink'
import './Profile.css'

const Profile = ({ user, onLogout }) => {
  const [profileData, setProfileData] = useState({
    firstName: '',
    lastName: '',
    phone: '',
    email: ''
  })
  const [passwordData, setPasswordData] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  })
  const [errors, setErrors] = useState({})
  const [passwordErrors, setPasswordErrors] = useState({})
  const [isLoading, setIsLoading] = useState(false)
  const [success, setSuccess] = useState('')

  useEffect(() => {
    if (user) {
      setProfileData({
        firstName: user.first_name || '',
        lastName: user.last_name || '',
        phone: user.phone || '',
        email: user.email || ''
      })
    }
  }, [user])

  const handleProfileChange = (e) => {
    const { name, value } = e.target
    setProfileData(prev => ({
      ...prev,
      [name]: value
    }))
    
    if (errors[name]) {
      setErrors(prev => ({
        ...prev,
        [name]: ''
      }))
    }
  }

  const validateProfile = () => {
    const newErrors = {}
    
    if (!profileData.firstName.trim()) {
      newErrors.firstName = 'Имя обязательно'
    }
    
    if (!profileData.lastName.trim()) {
      newErrors.lastName = 'Фамилия обязательна'
    }
    
    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleUpdateProfile = async () => {
    if (!validateProfile()) {
      return
    }
    
    setIsLoading(true)
    setSuccess('')
    
    try {
      const token = localStorage.getItem('token')
      
      const backendData = {
        first_name: profileData.firstName,
        last_name: profileData.lastName,
        phone: profileData.phone
      }
      
      const response = await fetch('/api/user/profile', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(backendData)
      })
      
      const data = await response.json()
      
      if (!response.ok) {
        throw new Error(data.error || 'Ошибка обновления профиля')
      }
      
      setSuccess('Профиль успешно обновлен')
      
      const updatedUser = { ...user, ...backendData }
      localStorage.setItem('user', JSON.stringify(updatedUser))
      
      setTimeout(() => setSuccess(''), 3000)
    } catch (error) {
      setErrors({ general: error.message })
    } finally {
      setIsLoading(false)
    }
  }

  const handlePasswordChange = (e) => {
    const { name, value } = e.target
    setPasswordData(prev => ({
      ...prev,
      [name]: value
    }))
    
    if (passwordErrors[name]) {
      setPasswordErrors(prev => ({
        ...prev,
        [name]: ''
      }))
    }
  }

  const validatePassword = () => {
    const newErrors = {}
    
    if (!passwordData.currentPassword) {
      newErrors.currentPassword = 'Текущий пароль обязателен'
    }
    
    if (!passwordData.newPassword) {
      newErrors.newPassword = 'Новый пароль обязателен'
    }
    
    if (passwordData.newPassword !== passwordData.confirmPassword) {
      newErrors.confirmPassword = 'Пароли не совпадают'
    }
    
    if (passwordData.newPassword.length < 6) {
      newErrors.newPassword = 'Пароль должен содержать минимум 6 символов'
    }
    
    setPasswordErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleChangePassword = async () => {
    if (!validatePassword()) {
      return
    }
    
    setIsLoading(true)
    setSuccess('')
    
    try {
      const token = localStorage.getItem('token')
      
      const response = await fetch('/api/user/password', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          current_password: passwordData.currentPassword,
          new_password: passwordData.newPassword
        })
      })
      
      const data = await response.json()
      
      if (!response.ok) {
        throw new Error(data.error || 'Ошибка смены пароля')
      }
      
      setSuccess('Пароль успешно изменен')
      setPasswordData({
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      })
      
      setTimeout(() => setSuccess(''), 3000)
    } catch (error) {
      setErrors({ general: error.message })
    } finally {
      setIsLoading(false)
    }
  }

  const handleDeleteAccount = async () => {
    if (!window.confirm('Вы уверены, что хотите удалить аккаунт? Это действие необратимо!')) {
      return
    }
    
    setIsLoading(true)
    
    try {
      const token = localStorage.getItem('token')
      
      const response = await fetch('/api/user/account', {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })
      
      if (!response.ok) {
        throw new Error('Ошибка удаления аккаунта')
      }
      
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      onLogout()
    } catch (error) {
      setErrors({ general: error.message })
    } finally {
      setIsLoading(false)
    }
  }

  if (!user) {
    return <div className="profile-container">Загрузка...</div>
  }

  return (
    <div className="profile-container">
      <div className="profile-header">
        <h1 className="profile-title">👤 Личный кабинет</h1>
        <button className="logout-button" onClick={onLogout}>
          🚪 Выйти
        </button>
      </div>

      {success && (
        <div className="success-message">
          ✅ {success}
        </div>
      )}

      {errors.general && (
        <div className="error-message">
          ❌ {errors.general}
        </div>
      )}

      <div className="profile-content">
        <div className="profile-section">
          <div className="section-header">
            <h2 className="section-title">📋 Информация о профиле</h2>
            <button className="edit-button">✏️ Редактировать</button>
          </div>

          <div className="profile-info">
            <div className="info-row">
              <div className="info-group">
                <label className="info-label">📧 Email</label>
                <input
                  type="email"
                  name="email"
                  className="info-input"
                  value={profileData.email}
                  onChange={handleProfileChange}
                  disabled={true}
                  title="Email нельзя изменить"
                />
              </div>
              <div className="info-group">
                <label className="info-label">📱 Телефон</label>
                <input
                  type="tel"
                  name="phone"
                  className="info-input"
                  value={profileData.phone}
                  onChange={handleProfileChange}
                  placeholder="+7 (999) 123-45-67"
                />
              </div>
            </div>

            <div className="info-row">
              <div className="info-group">
                <label className="info-label">👤 Имя</label>
                <input
                  type="text"
                  name="firstName"
                  className="info-input"
                  value={profileData.firstName}
                  onChange={handleProfileChange}
                  disabled={true}
                />
              </div>
              <div className="info-group">
                <label className="info-label">👥 Фамилия</label>
                <input
                  type="text"
                  name="lastName"
                  className="info-input"
                  value={profileData.lastName}
                  onChange={handleProfileChange}
                  disabled={true}
                />
              </div>
            </div>
          </div>
        </div>

        <div className="profile-section">
          <div className="section-header">
            <h2 className="section-title">🔒 Безопасность</h2>
          </div>

          <div className="password-form">
            <div className="info-group">
              <label className="info-label">🔑 Текущий пароль</label>
              <input
                type="password"
                name="currentPassword"
                className="info-input"
                value={passwordData.currentPassword}
                onChange={handlePasswordChange}
                placeholder="Введите текущий пароль"
              />
              {passwordErrors.currentPassword && (
                <div className="error-text">{passwordErrors.currentPassword}</div>
              )}
            </div>

            <div className="info-group">
              <label className="info-label">🆕 Новый пароль</label>
              <input
                type="password"
                name="newPassword"
                className="info-input"
                value={passwordData.newPassword}
                onChange={handlePasswordChange}
                placeholder="Минимум 6 символов"
              />
              {passwordErrors.newPassword && (
                <div className="error-text">{passwordErrors.newPassword}</div>
              )}
            </div>

            <div className="info-group">
              <label className="info-label">🔄 Подтверждение нового пароля</label>
              <input
                type="password"
                name="confirmPassword"
                className="info-input"
                value={passwordData.confirmPassword}
                onChange={handlePasswordChange}
                placeholder="Повторите новый пароль"
              />
              {passwordErrors.confirmPassword && (
                <div className="error-text">{passwordErrors.confirmPassword}</div>
              )}
            </div>

            <button className="save-button" onClick={handleChangePassword} disabled={isLoading}>
              🔐 Сменить пароль
            </button>
          </div>
        </div>

        <div className="profile-section danger-zone">
          <div className="section-header">
            <h2 className="section-title">⚠️ Опасная зона</h2>
          </div>

          <p className="danger-text">
            🚨 <strong>Внимание!</strong> Удаление аккаунта является необратимым действием. 
            Все ваши данные, включая историю заказов и личную информацию, будут безвозвратно удалены.
          </p>

          <button className="delete-button" onClick={handleDeleteAccount} disabled={isLoading}>
            🗑️ Удалить аккаунт навсегда
          </button>
        </div>

        <OAuthLink 
          user={user} 
          onLinkSuccess={() => setSuccess('OAuth аккаунт успешно привязан!')}
        />
      </div>
    </div>
  )
}

export default Profile
