import React, { useState, useEffect } from 'react'
import Header from './components/Header'
import ProductGrid from './components/ProductGrid'
import FeaturedProducts from './components/FeaturedProducts'
import CategoryFilter from './components/CategoryFilter'
import SearchBar from './components/SearchBar'
import './App.css'

function App() {
  const [products, setProducts] = useState([])
  const [featuredProducts, setFeaturedProducts] = useState([])
  const [categories, setCategories] = useState([])
  const [selectedCategory, setSelectedCategory] = useState('')
  const [searchQuery, setSearchQuery] = useState('')
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const limit = 12

  useEffect(() => {
    fetchCategories()
    fetchFeaturedProducts()
  }, [])

  useEffect(() => {
    if (searchQuery) {
      handleSearch(searchQuery)
    } else if (selectedCategory) {
      fetchProductsByCategory(selectedCategory)
    } else {
      fetchProducts()
    }
  }, [page, selectedCategory, searchQuery])

  const fetchProducts = async () => {
    setLoading(true)
    try {
      const response = await fetch(`/api/products?page=${page}&limit=${limit}`)
      const data = await response.json()
      setProducts(data.products || [])
      setTotal(data.total || 0)
    } catch (error) {
      console.error('Error fetching products:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchProductsByCategory = async (category) => {
    setLoading(true)
    try {
      const response = await fetch(`/api/products/category/${encodeURIComponent(category)}?page=${page}&limit=${limit}`)
      const data = await response.json()
      setProducts(data.products || [])
      setTotal(data.total || 0)
    } catch (error) {
      console.error('Error fetching products by category:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchFeaturedProducts = async () => {
    try {
      const response = await fetch('/api/products/featured?limit=8')
      const data = await response.json()
      setFeaturedProducts(data.products || [])
    } catch (error) {
      console.error('Error fetching featured products:', error)
    }
  }

  const fetchCategories = async () => {
    try {
      const response = await fetch('/api/categories')
      const data = await response.json()
      setCategories(data.categories || [])
    } catch (error) {
      console.error('Error fetching categories:', error)
    }
  }

  const handleSearch = async (query) => {
    if (!query.trim()) {
      setSearchQuery('')
      setSelectedCategory('')
      setPage(1)
      fetchProducts()
      return
    }

    setLoading(true)
    setSearchQuery(query)
    setSelectedCategory('')
    setPage(1)
    try {
      const response = await fetch(`/api/products/search?q=${encodeURIComponent(query)}&page=${page}&limit=${limit}`)
      const data = await response.json()
      setProducts(data.products || [])
      setTotal(data.total || 0)
    } catch (error) {
      console.error('Error searching products:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleCategorySelect = (category) => {
    setSelectedCategory(category)
    setSearchQuery('')
    setPage(1)
  }

  const handleClearFilters = () => {
    setSelectedCategory('')
    setSearchQuery('')
    setPage(1)
  }

  const totalPages = Math.ceil(total / limit)

  return (
    <div className="app">
      <Header />
      <main className="main-content">
        <div className="container">
          <h1 className="main-title">Добро пожаловать в Tenderness</h1>
          <p className="main-subtitle">Премиум товары для особых моментов</p>

          <SearchBar onSearch={handleSearch} />

          <CategoryFilter
            categories={categories}
            selectedCategory={selectedCategory}
            onSelectCategory={handleCategorySelect}
            onClearFilters={handleClearFilters}
          />

          {!selectedCategory && !searchQuery && featuredProducts.length > 0 && (
            <section className="featured-section">
              <h2 className="section-title">Рекомендуемые товары</h2>
              <FeaturedProducts products={featuredProducts} />
            </section>
          )}

          <section className="products-section">
            <div className="section-header">
              <h2 className="section-title">
                {searchQuery ? `Результаты поиска: "${searchQuery}"` : 
                 selectedCategory ? `Категория: ${selectedCategory}` : 
                 'Все товары'}
              </h2>
              {total > 0 && (
                <span className="products-count">Найдено: {total}</span>
              )}
            </div>

            {loading ? (
              <div className="loading">Загрузка...</div>
            ) : products.length > 0 ? (
              <>
                <ProductGrid products={products} />
                {totalPages > 1 && (
                  <div className="pagination">
                    <button
                      onClick={() => setPage(p => Math.max(1, p - 1))}
                      disabled={page === 1}
                      className="pagination-btn"
                    >
                      Назад
                    </button>
                    <span className="pagination-info">
                      Страница {page} из {totalPages}
                    </span>
                    <button
                      onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                      disabled={page === totalPages}
                      className="pagination-btn"
                    >
                      Вперед
                    </button>
                  </div>
                )}
              </>
            ) : (
              <div className="no-products">
                <p>Товары не найдены</p>
              </div>
            )}
          </section>
        </div>
      </main>
    </div>
  )
}

export default App

