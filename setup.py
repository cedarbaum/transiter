from setuptools import setup, find_packages

metadata = {}
with open("transiter/__metadata__.py") as f:
    exec(f.read(), metadata)
version = metadata["__version__"]

setup(
    name="transiter",
    version=version,
    author="James Fennell",
    author_email="jamespfennell@gmail.com",
    description="HTTP web service for transit data",
    url="https://github.com/jamespfennell/transiter",
    packages=find_packages(),
    license="MIT",
    entry_points={"console_scripts": ["transiterclt = transiter.clt:transiter_clt"]},
    install_requires=[
        "alembic==1.4.2",
        "apscheduler==3.6.3",
        "celery==4.4.2",
        "click==7.1.1",
        "decorator==4.4.2",
        "flask==1.1.2",
        "gunicorn==20.0.4",
        "inflection==0.4.0",
        "Jinja2==2.11.2",
        "protobuf==3.11.3",
        "psycopg2-binary==2.8.5",
        "pytimeparse==1.1.8",
        "pytz==2019.3",
        "requests==2.23.0",
        "sqlalchemy==1.3.16",
        "strictyaml==1.0.6",
    ],
)
