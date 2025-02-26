import re
import os
from bs4 import BeautifulSoup
import sqlalchemy
from sqlalchemy.orm import DeclarativeBase, Session
import requests


URL = 'https://diablo2.io/runewords/'
PATH = 'rune_image'



class Base(DeclarativeBase):
    pass


class Runes(Base):
    __tablename__ = 'rs'

    id = sqlalchemy.Column(sqlalchemy.INTEGER, primary_key=True)
    name = sqlalchemy.Column(sqlalchemy.String(25), unique=True, nullable=False)
    image = sqlalchemy.Column(sqlalchemy.BLOB, nullable=False)


class RuneWords(Base):
    __tablename__ = 'rws'

    id = sqlalchemy.Column(sqlalchemy.INTEGER, primary_key=True)
    name = sqlalchemy.Column(sqlalchemy.String(50), unique=True, nullable=False)
    item_class = sqlalchemy.Column(sqlalchemy.String(50), nullable=False)
    ladder = sqlalchemy.Column(sqlalchemy.Boolean, nullable=False)
    sockets = sqlalchemy.Column(sqlalchemy.INTEGER, nullable=False)
    lvl = sqlalchemy.Column(sqlalchemy.INTEGER, nullable=False)
    runes = sqlalchemy.Column(sqlalchemy.String(50), nullable=False)
    description = sqlalchemy.Column(sqlalchemy.String(255), nullable=False)


def extract_number(filename):
    match = re.search(r'\d+', filename)
    return int(match.group()) if match else 0


def collect_rw(url: str) -> list:
    """
    Parse url, collect data.
    return data list
    """
    response = requests.get(url).text
    page = BeautifulSoup(response, 'html.parser')

    runewords = page.find_all('article')
    runeword = []

    for row in runewords:
        temp = {}
        temp_runes = []
        temp_item_class = []
        
        temp['name'] = row.find('h3', class_='z-sort-name').text

        for rune in row.find_all('span', class_='ajax_catch'):
            temp_runes.append(rune.find('div').get('data-background-image').split('rune')[-1].split('_')[0])
        
        temp['runes'] = ', '.join(temp_runes)
        temp['lvl'] = int(row.find('span', class_='zso_rwlvlrq').text)
        temp['sockets'] = int(row.find('span', class_='zso_rwsock').text)
        
        for item in row.find_all('div', class_='z-vf-hide'):
            for iclass in item.find_all('a'):
                cleaned_text = re.sub(r'[\n\r\t]', '', iclass.text.strip())
                temp_item_class.append(cleaned_text.split(' ')[-1])

            for description in item.find_all('span', class_='z-smallstats'):
                if not description:
                    continue
                print(1, description, description.text, sep='\n')
                cleaned_text = re.sub(r'[\n\r\t]', '', description.text.strip())
                temp['description'] = cleaned_text
                print(cleaned_text)
            # return
        temp['item_class'] = ', '.join(temp_item_class)

        if row.find('span', class_='zi zi-bb zi-ladder z-ic'):
            temp['ladder'] = True
        else:
            temp['ladder'] = False

        runeword.append(temp)

    return runeword


def collect_r(path: str) -> list:
    """
    Read images
    """
    runes = []
    rune_images = os.listdir(path)
    sorted_rune_image = sorted(rune_images, key=extract_number)

    for image in sorted_rune_image:
        temp = {}
        id_, name, _  = image.split('.')
        with open('\\'.join((path, image)), 'rb') as f:
            temp['id'] = id_
            temp['name'] = name
            temp['image'] = f.read()
        runes.append(temp)
    return runes


def main():
    runewords = collect_rw(url=URL)
    # runes = collect_r(path=PATH)
    # engine = sqlalchemy.create_engine(r'sqlite:///..\sqlite.db', echo=True)
    # Base.metadata.create_all(bind=engine)
    
    # for runeword in runewords:
    #     with Session(autoflush=False, bind=engine) as db:
    #         rw = RuneWords(
    #             name=runeword['name'],
    #             item_class=runeword['item_class'],
    #             ladder = runeword['ladder'],
    #             sockets = runeword['sockets'],
    #             lvl = runeword['lvl'],
    #             runes = runeword['runes'],
    #             description = runeword['description']
    #             )
    #         db.add(rw)
    #         db.commit()

    # for rune in runes:
    #     with Session(autoflush=False, bind=engine) as db:
    #         r = Runes(name=rune['name'], image=rune['image'])
    #         db.add(r)
    #         db.commit()
    

if __name__ == '__main__':
    main()
