import React from 'react'
import './ProductCard.css'

function ProductCard({ product }) {
  return (
    <div className="product-card">
      <div className="product-image-container">
        <img
          src={product.image_url || 'https://via.placeholder.com/400'}
          alt={product.name}
          className="product-image"
          onError={(e) => {
            e.target.src = 'https://via.placeholder.com/400'
          }}
        />
        {!product.in_stock && (
          <div className="out-of-stock-badge">Нет в наличии</div>
        )}
        {product.rating > 0 && (
          <div className="rating-badge">
            ⭐ {product.rating.toFixed(1)}
          </div>
        )}
      </div>
      <div className="product-info">
        <h3 className="product-name">{product.name}</h3>
        <p className="product-category">{product.category}</p>
        {product.description && (
          <p className="product-description">
            {product.description.length > 100
              ? `${product.description.substring(0, 100)}...`
              : product.description}
          </p>
        )}
        <div className="product-footer">
          <span className="product-price">
            {product.price.toLocaleString('ru-RU')} ₽
          </span>
          <button className="add-to-cart-btn">В корзину</button>
        </div>
      </div>
    </div>
  )
}

export default ProductCard

