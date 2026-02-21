import React from 'react'
import './OAuthButtons.css'

const OAuthButtons = ({ onOAuthLogin }) => {
  const handleGoogleLogin = () => {
    window.location.href = '/api/oauth2/google/auth'
  }

  const handleGitHubLogin = () => {
    window.location.href = '/api/oauth2/github/auth'
  }

  return (
    <div className="oauth-buttons">
      <div className="oauth-divider">
        <span className="oauth-divider-text">или войдите через</span>
      </div>
      
      <button 
        className="oauth-button google-button"
        onClick={handleGoogleLogin}
        type="button"
      >
        <span className="oauth-icon">🔍</span>
        <span className="oauth-text">Google</span>
      </button>
      
      <button 
        className="oauth-button github-button"
        onClick={handleGitHubLogin}
        type="button"
      >
        <span className="oauth-icon">🐙</span>
        <span className="oauth-text">GitHub</span>
      </button>
    </div>
  )
}

export default OAuthButtons
