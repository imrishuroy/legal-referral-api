CREATE TABLE firms (
    firm_id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    logo_url VARCHAR NOT NULL,
    org_type VARCHAR NOT NULL,
    website VARCHAR NOT NULL,
    location VARCHAR NOT NULL,
    about TEXT NOT NULL
);

INSERT INTO firms (name, logo_url, org_type, website, location, about)
VALUES
    ('Sullivan & Cromwell LLP', 'https://legal-referral.s3.ap-south-1.amazonaws.com/Sullivan.webp', 'Law Firm', 'www.sullcrom.com', 'New York', 'One of the most prestigious law firms, Sullivan & Cromwell subscribes to a generalist approach, allowing lawyers to work across industries and subgroups. Sullivan & Cromwell boasts an army of over 800 attorneys across eight countries.'),
    ('Wachtell, Lipton, Rosen & Katz', 'https://legal-referral.s3.ap-south-1.amazonaws.com/watch.jpeg', 'Law Firm', 'www.wlrk.com', 'New York', 'M&A giant and a leader in the New York legal market, Wachtell Lipton is one of the most elite law firms in the industry. Known for its high-profile matters and above-market compensation, Wachtell Lipton is home to a passionate group of lawyers who embrace hard work and professionalism. The firm is small by elite, BigLaw standards, which fosters a collegial atmosphere, and has spurned domestic or international expansion by maintaining a sole office in NYC.'),
    ('Cravath, Swaine & Moore LLP', 'https://legal-referral.s3.ap-south-1.amazonaws.com/Sullivan.webp', 'Law Firm', 'www.cravath.com', 'New York', 'Cravath, Swaine & Moore is one of the most prestigious law firms in the world. The firm is known for its high-profile clients and high-stakes matters, and it is considered one of the most elite law firms in the United States. Cravath is a general practice firm, but it is best known for its corporate work. The firm is also known for its collegial atmosphere and its commitment to diversity and pro bono work.'),
    ('Davis Polk & Wardwell LLP', 'https://legal-referral.s3.ap-south-1.amazonaws.com/davis_polk__wardwell_llp_logo.jpeg', 'Law Firm', 'www.davispolk.com', 'New York', 'Davis Polk & Wardwell is one of the most prestigious law firms in the world. The firm is known for its high-profile clients and high-stakes matters, and it is considered one of the most elite law firms in the United States. Davis Polk is a general practice firm, but it is best known for its corporate work. The firm is also known for its collegial atmosphere and its commitment to diversity and pro bono work.');
