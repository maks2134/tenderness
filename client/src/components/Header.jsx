import React from 'react'
import './Header.css'

function Header() {
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
        </div>
      </div>
    </header>
  )
}

export default Header

