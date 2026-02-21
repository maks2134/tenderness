import React from 'react'
import './Header.css'

function Header({ user, onLogin, onLogout, onProfile }) {
  return (
    <header className="header">
      <div className="container">
        <div className="header-content">
          <h1 className="logo">Tenderness</h1>
          <nav className="nav">
            <a href="#" className="nav-link">Главная</a>
            <a href="#" className="nav-link">Каталог</a>
            <a href="#" className="nav-link">О нас</a>
            <a href="#" className="nav-link">Контакты</a>
          </nav>
          <div className="auth-section">
            {user ? (
              <div className="user-menu">
                <span className="welcome-text">Добро пожаловать, {user.first_name}!</span>
                <button className="auth-button profile-button" onClick={onProfile}>
                  Личный кабинет
                </button>
                <button className="auth-button logout-button" onClick={onLogout}>
                  Выйти
                </button>
              </div>
            ) : (
              <div className="auth-buttons">
                <button className="auth-button login-button" onClick={onLogin}>
                  Войти
                </button>
                <button className="auth-button register-button" onClick={onLogin}>
                  Регистрация
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  )
}

export default Header

