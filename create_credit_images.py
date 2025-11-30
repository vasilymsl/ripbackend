#!/usr/bin/env python3
"""Создаем простые SVG картинки для кредитных продуктов"""

credits = [
    (1, 'Потребительский кредит', '1e3a5f', 'percent'),
    (2, 'Кредит наличными', '2c5aa0', 'wallet'),
    (3, 'Кредит под залог', '3a6fb0', 'bank'),
    (4, 'Кредитная карта', '4a7fc1', 'card'),
    (5, 'Автокредит', '5a8fd1', 'car'),
    (6, 'Ипотека', '6a9fe2', 'home'),
]

sql_updates = []

for credit_id, title, color, icon in credits:
    # Создаем простой SVG
    svg = f'''<svg xmlns="http://www.w3.org/2000/svg" width="400" height="300" viewBox="0 0 400 300">
  <rect width="400" height="300" fill="#{color}"/>
  <circle cx="200" cy="120" r="50" fill="white" opacity="0.2"/>
  <text x="200" y="180" font-family="Arial, sans-serif" font-size="22" fill="white" text-anchor="middle" font-weight="bold">{title}</text>
  <text x="200" y="220" font-family="Arial, sans-serif" font-size="16" fill="white" text-anchor="middle" opacity="0.8">Кредитный продукт</text>
</svg>'''
    
    # Кодируем в base64
    import base64
    svg_base64 = base64.b64encode(svg.encode('utf-8')).decode('utf-8')
    data_uri = f"data:image/svg+xml;base64,{svg_base64}"
    
    sql_updates.append(f"UPDATE credits SET image_url = '{data_uri}' WHERE id = {credit_id};")

print('\n'.join(sql_updates))

