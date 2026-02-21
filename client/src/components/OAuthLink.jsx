import React, { useState } from 'react'
import './OAuthLink.css'

const OAuthLink = ({ user, onLinkSuccess }) => {
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')

  const handleGoogleLink = () => {
    setIsLoading(true)
    setError('')
    
    const token = localStorage.getItem('token')
    if (!token) {
      setError('Сначала войдите в аккаунт')
      setIsLoading(false)
      return
    }

    // Redirect to Google OAuth
    window.location.href = `/api/oauth2/google/auth?state=${encodeURIComponent(token)}`
  }

  const handleGitHubLink = () => {
    setIsLoading(true)
    setError('')
    
    const token = localStorage.getItem('token')
    if (!token) {
      setError('Сначала войдите в аккаунт')
      setIsLoading(false)
      return
    }

    // Redirect to GitHub OAuth
    window.location.href = `/api/oauth2/github/auth?state=${encodeURIComponent(token)}`
  }

  const isLinked = (provider) => {
    switch (provider) {
      case 'google':
        return user.google_id
      case 'github':
        return user.github_id
      default:
        return false
    }
  }

  return (
    <div className="oauth-link">
      <h3 className="oauth-link-title">Привязка аккаунтов</h3>
      <p className="oauth-link-description">
        Привяжите ваши аккаунты Google и GitHub для быстрого входа в будущем
      </p>

      {error && (
        <div className="oauth-error">
          ❌ {error}
        </div>
      )}

      <div className="oauth-providers">
        <div className="oauth-provider">
          <div className="provider-info">
            <span className="provider-icon">🔍</span>
            <span className="provider-name">Google</span>
            {isLinked('google') ? (
              <span className="linked-status">✅ Привязан</span>
            ) : (
              <button 
                className="link-button"
                onClick={handleGoogleLink}
                disabled={isLoading}
              >
                {isLoading ? 'Привязка...' : 'Привязать Google'}
              </button>
            )}
          </div>
        </div>

        <div className="oauth-provider">
          <div className="provider-info">
            <span className="provider-icon">🐙</span>
            <span className="provider-name">GitHub</span>
            {isLinked('github') ? (
              <span className="linked-status">✅ Привязан</span>
            ) : (
              <button 
                className="link-button"
                onClick={handleGitHubLink}
                disabled={isLoading}
              >
                {isLoading ? 'Привязка...' : 'Привязать GitHub'}
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default OAuthLink
