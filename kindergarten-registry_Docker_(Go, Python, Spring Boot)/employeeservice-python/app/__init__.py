from flask import Flask
from flask_cors import CORS
import os

def create_app():
    app = Flask(__name__)
    CORS(app)
    
    # APM Configuration - Only if APM is enabled
    if os.getenv('APM_ENABLED', 'False').lower() == 'true':
        from elasticapm.contrib.flask import ElasticAPM
        app.config['ELASTIC_APM'] = {
            'SERVICE_NAME': 'employee-service-python',
            'SECRET_TOKEN': os.getenv('APM_SECRET_TOKEN', ''),
            'SERVER_URL': os.getenv('APM_SERVER_URL', ''),
            'ENVIRONMENT': os.getenv('ENVIRONMENT', 'development'),
            'ENABLED': True
        }
        ElasticAPM(app)
        print("APM enabled for employee service")
    else:
        print("APM disabled for employee service")
    
    # Import and register blueprints
    from .routes import employee_bp
    app.register_blueprint(employee_bp)
    
    return app