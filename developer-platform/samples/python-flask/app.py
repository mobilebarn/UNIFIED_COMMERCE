# Sample Python Flask application for Retail OS
import os
import requests
from flask import Flask, jsonify, render_template_string
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

app = Flask(__name__)

# Retail OS API client
class RetailOSClient:
    def __init__(self):
        self.api_endpoint = os.getenv('API_ENDPOINT', 'http://localhost:4000/graphql')
        self.client_id = os.getenv('CLIENT_ID')
        self.client_secret = os.getenv('CLIENT_SECRET')
        self.access_token = None

    def authenticate(self):
        # In a real implementation, you would authenticate with the API
        # For this sample, we'll just return a mock token
        self.access_token = 'mock-access-token'
        return self.access_token

    def make_request(self, query, variables=None):
        if not self.access_token:
            self.authenticate()

        headers = {
            'Authorization': f'Bearer {self.access_token}',
            'Content-Type': 'application/json'
        }

        payload = {
            'query': query,
            'variables': variables or {}
        }

        try:
            response = requests.post(self.api_endpoint, json=payload, headers=headers)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            print(f'API request failed: {e}')
            raise

    # Product methods
    def get_products(self, limit=10):
        query = '''
            query GetProducts($limit: Int) {
                products(filter: { limit: $limit }) {
                    id
                    title
                    handle
                    description
                    featuredImage
                    priceRange {
                        minVariantPrice
                    }
                    tags
                }
            }
        '''

        response = self.make_request(query, {'limit': limit})
        return response['data']['products']

    def get_product_by_id(self, id):
        query = '''
            query GetProduct($id: ID!) {
                product(id: $id) {
                    id
                    title
                    handle
                    description
                    featuredImage
                    images {
                        src
                        altText
                    }
                    priceRange {
                        minVariantPrice
                    }
                    variants {
                        id
                        title
                        price
                        inventoryQuantity
                    }
                    tags
                }
            }
        '''

        response = self.make_request(query, {'id': id})
        return response['data']['product']

    # User methods
    def get_current_user(self):
        query = '''
            query {
                currentUser {
                    id
                    email
                    firstName
                    lastName
                }
            }
        '''

        response = self.make_request(query)
        return response['data']['currentUser']

    # Order methods
    def get_orders(self, limit=10):
        query = '''
            query GetOrders($limit: Int) {
                orders(filter: { limit: $limit }) {
                    id
                    orderNumber
                    status
                    totalPrice
                    createdAt
                    lineItems {
                        id
                        name
                        quantity
                        price
                    }
                }
            }
        '''

        response = self.make_request(query, {'limit': limit})
        return response['data']['orders']

# Initialize client
client = RetailOSClient()

@app.route('/')
def index():
    return render_template_string('''
    <!DOCTYPE html>
    <html>
    <head>
        <title>Retail OS Python Sample</title>
        <style>
            body { font-family: Arial, sans-serif; margin: 40px; }
            .endpoint { margin: 10px 0; padding: 10px; background: #f5f5f5; }
        </style>
    </head>
    <body>
        <h1>Retail OS Python Sample</h1>
        <p>This sample demonstrates how to integrate with the Retail OS API using Python and Flask.</p>
        
        <h2>API Endpoints</h2>
        <div class="endpoint"><a href="/api/products">/api/products</a> - Get products</div>
        <div class="endpoint"><a href="/api/user">/api/user</a> - Get current user</div>
        <div class="endpoint"><a href="/api/orders">/api/orders</a> - Get orders</div>
        
        <h2>Documentation</h2>
        <p>See <a href="https://docs.retail-os.com">API Documentation</a> for more details.</p>
    </body>
    </html>
    ''')

@app.route('/api/products')
def get_products():
    try:
        limit = int(request.args.get('limit', 10))
        products = client.get_products(limit)
        return jsonify(products)
    except Exception as e:
        return jsonify({'error': 'Failed to fetch products'}), 500

@app.route('/api/products/<id>')
def get_product(id):
    try:
        product = client.get_product_by_id(id)
        return jsonify(product)
    except Exception as e:
        return jsonify({'error': 'Failed to fetch product'}), 500

@app.route('/api/user')
def get_user():
    try:
        user = client.get_current_user()
        return jsonify(user)
    except Exception as e:
        return jsonify({'error': 'Failed to fetch user'}), 500

@app.route('/api/orders')
def get_orders():
    try:
        limit = int(request.args.get('limit', 10))
        orders = client.get_orders(limit)
        return jsonify(orders)
    except Exception as e:
        return jsonify({'error': 'Failed to fetch orders'}), 500

if __name__ == '__main__':
    print(f"Retail OS sample app starting...")
    print(f"API endpoint: {client.api_endpoint}")
    app.run(host='0.0.0.0', port=int(os.getenv('PORT', 3000)), debug=True)