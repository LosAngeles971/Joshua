docker run \
        --rm \
        -v $CI_PROJECT_DIR/whitepaper/:/documents/ \
        asciidoctor/docker-asciidoctor:latest \
        asciidoctor -r asciidoctor-pdf -r asciidoctor-mathematical -b pdf \
        /documents/0.cover.adoc \
        -o /documents/joshua.pdf