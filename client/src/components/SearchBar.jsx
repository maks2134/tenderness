import React, { useState } from 'react'
import './SearchBar.css'

function SearchBar({ onSearch }) {
  const [query, setQuery] = useState('')

  const handleSubmit = (e) => {
    e.preventDefault()
    onSearch(query)
  }

  const handleChange = (e) => {
    setQuery(e.target.value)
    if (!e.target.value.trim()) {
      onSearch('')
    }
  }

  return (
    <form className="search-bar" onSubmit={handleSubmit}>
      <input
        type="text"
        className="search-input"
        placeholder="Поиск товаров..."
        value={query}
        onChange={handleChange}
      />
      <button type="submit" className="search-btn">
        Поиск
      </button>
    </form>
  )
}

export default SearchBar

