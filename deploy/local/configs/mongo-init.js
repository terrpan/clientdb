db = db.getSiblingDB('clientdb');

db.createCollection('clients');
db.createCollection('contacts');
db.createCollection('services');
db.clients.insertMany(
  [
    {
      "client_name": "Larkin, Blick and Cremin",
      "slack_channel": "Subin",
      "web_url": "https://wp.com"
    },
    {
      "client_name": "Legros-Hayes",
      "slack_channel": "Wrapsafe",
      "web_url": "http://go.com"
    },
    {
      "client_name": "Okuneva, Heaney and Hansen",
      "slack_channel": "Flowdesk",
      "web_url": "http://webmd.com"
    },
    {
      "client_name": "Mante, Dach and Crona",
      "slack_channel": "Overhold",
      "web_url": "https://cocolog-nifty.com"
    },
    {
      "client_name": "Rolfson Inc",
      "slack_channel": "Lotstring",
      "web_url": "http://php.net"
    },
    {
      "client_name": "McKenzie-Aufderhar",
      "slack_channel": "Sonair",
      "web_url": "http://linkedin.com"
    },
    {
      "client_name": "Glover, Morar and Bins",
      "slack_channel": "Matsoft",
      "web_url": "https://tamu.edu"
    },
    {
      "client_name": "Schmitt-Kuhic",
      "slack_channel": "Cookley",
      "web_url": "https://arizona.edu"
    },
    {
      "client_name": "Bernhard-Heller",
      "slack_channel": "Otcom",
      "web_url": "https://mediafire.com"
    },
    {
      "client_name": "Schamberger LLC",
      "slack_channel": "Job",
      "web_url": "http://msn.com"
    }
  ]
);

db.contacts.insertMany(
  [
    {
      "first_name": "Leah",
      "last_name": "Franz-Schoninger",
      "email": "lfranzschoninger0@indiatimes.com",
      "role": "Construction Expeditor",
      "phone_number": "969-986-3311"
    },
    {
      "first_name": "Carroll",
      "last_name": "Harness",
      "email": "charness1@google.co.jp",
      "role": "Subcontractor",
      "phone_number": "864-941-5264"
    },
    {
      "first_name": "Janka",
      "last_name": "Willbourne",
      "email": "jwillbourne2@merriam-webster.com",
      "role": "Project Manager",
      "phone_number": "711-333-7020"
    },
    {
      "first_name": "Jeromy",
      "last_name": "Becket",
      "email": "jbecket3@vinaora.com",
      "role": "Engineer",
      "phone_number": "220-506-6329"
    },
    {
      "first_name": "Mendy",
      "last_name": "Dorre",
      "email": "mdorre4@wikia.com",
      "role": "Subcontractor",
      "phone_number": "648-887-5402"
    },
    {
      "first_name": "Roda",
      "last_name": "Torvey",
      "email": "rtorvey5@epa.gov",
      "role": "Construction Foreman",
      "phone_number": "525-807-9796"
    },
    {
      "first_name": "Berti",
      "last_name": "Gayne",
      "email": "bgayne6@etsy.com",
      "role": "Construction Foreman",
      "phone_number": "853-126-7362"
    },
    {
      "first_name": "Eziechiele",
      "last_name": "Booler",
      "email": "ebooler7@constantcontact.com",
      "role": "Construction Expeditor",
      "phone_number": "543-867-1979"
    },
    {
      "first_name": "Nichols",
      "last_name": "Jerrard",
      "email": "njerrard8@stumbleupon.com",
      "role": "Subcontractor",
      "phone_number": "905-693-6375"
    },
    {
      "first_name": "Amery",
      "last_name": "Wistance",
      "email": "awistance9@patch.com",
      "role": "Subcontractor",
      "phone_number": "514-462-9967"
    }
  ]
);

db.services.insertMany(
  [
    {
      "service_name": "Redhold",
      "service_type": "br.com.uol.Lotstring",
      "service_owner": "Emmit Gearty",
      "service_description": "Dilation of Bladder Neck, Open Approach",
      "service_status": "Tools",
      "invoice_frequency": "Yearly",
      "Invoice_amount": 256,
      "management_fee": 90
    },
    {
      "service_name": "Duobam",
      "service_type": "com.ask.Subin",
      "service_owner": "Conan Lewsam",
      "service_description": "Destruction of Bilateral Seminal Vesicles, Open Approach",
      "service_status": "Jewelry",
      "invoice_frequency": "Monthly",
      "Invoice_amount": 598,
      "management_fee": 81
    },
    {
      "service_name": "Cardguard",
      "service_type": "com.studiopress.Latlux",
      "service_owner": "Laverna Spender",
      "service_description": "Destruction of Ileocecal Valve, Perc Endo Approach",
      "service_status": "Sports",
      "invoice_frequency": "Daily",
      "Invoice_amount": 417,
      "management_fee": 53
    },
    {
      "service_name": "Treeflex",
      "service_type": "com.godaddy.Sub-Ex",
      "service_owner": "Seymour Jenyns",
      "service_description": "Restrict Sigmoid Colon w Extralum Dev, Perc Endo",
      "service_status": "Sports",
      "invoice_frequency": "Weekly",
      "Invoice_amount": 196,
      "management_fee": 61
    },
    {
      "service_name": "Treeflex",
      "service_type": "us.imageshack.Wrapsafe",
      "service_owner": "Corene Lipscombe",
      "service_description": "Release Left Innominate Vein, Percutaneous Approach",
      "service_status": "Jewelry",
      "invoice_frequency": "Seldom",
      "Invoice_amount": 716,
      "management_fee": 93
    },
    {
      "service_name": "Tempsoft",
      "service_type": "com.dropbox.Kanlam",
      "service_owner": "Hewett Heatly",
      "service_description": "Insertion of Ext Fix into Skull, Perc Approach",
      "service_status": "Jewelry",
      "invoice_frequency": "Daily",
      "Invoice_amount": 287,
      "management_fee": 34
    },
    {
      "service_name": "Redhold",
      "service_type": "com.surveymonkey.Sonsing",
      "service_owner": "Donall Bontine",
      "service_description": "Extirpation of Matter from L Verteb Art, Open Approach",
      "service_status": "Kids",
      "invoice_frequency": "Daily",
      "Invoice_amount": 927,
      "management_fee": 60
    },
    {
      "service_name": "Aerified",
      "service_type": "cn.com.sina.Trippledex",
      "service_owner": "Heath Peddel",
      "service_description": "Dilation of Left Hand Artery with 3 Drug-elut, Perc Approach",
      "service_status": "Jewelry",
      "invoice_frequency": "Weekly",
      "Invoice_amount": 839,
      "management_fee": 56
    },
    {
      "service_name": "Viva",
      "service_type": "edu.princeton.Aerified",
      "service_owner": "Karie Titford",
      "service_description": "Bypass R Com Iliac Art to B Ext Ilia w Autol Vn, Perc Endo",
      "service_status": "Games",
      "invoice_frequency": "Seldom",
      "Invoice_amount": 76,
      "management_fee": 95
    },
    {
      "service_name": "Flowdesk",
      "service_type": "edu.stanford.Pannier",
      "service_owner": "Stan Wallenger",
      "service_description": "Plaque Radiation of Maxilla",
      "service_status": "Toys",
      "invoice_frequency": "Monthly",
      "Invoice_amount": 592,
      "management_fee": 72
    }
  ]
);