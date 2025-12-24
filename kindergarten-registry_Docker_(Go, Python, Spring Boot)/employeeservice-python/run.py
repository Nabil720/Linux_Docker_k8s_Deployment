import os
from app import create_app

# Initialize the app
app = create_app()

if __name__ == '__main__':
    port = int(os.getenv('PORT', 5003))
    debug = os.getenv('FLASK_ENV', 'production') == 'development'
    print(f"Starting Employee Service on port {port}")
    print(f"Debug mode: {debug}")
    app.run(host='0.0.0.0', port=port, debug=debug)