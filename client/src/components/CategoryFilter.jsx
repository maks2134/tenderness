import React from 'react'
import './CategoryFilter.css'

function CategoryFilter({ categories, selectedCategory, onSelectCategory, onClearFilters }) {
  if (categories.length === 0) return null

  return (
    <div className="category-filter">
      <div className="category-buttons">
        <button
          className={`category-btn ${!selectedCategory ? 'active' : ''}`}
          onClick={onClearFilters}
        >
          Все
        </button>
        {categories.map((category) => (
          <button
            key={category.id}
            className={`category-btn ${selectedCategory === category.name ? 'active' : ''}`}
            onClick={() => onSelectCategory(category.name)}
          >
            {category.name}
          </button>
        ))}
      </div>
    </div>
  )
}

export default CategoryFilter

