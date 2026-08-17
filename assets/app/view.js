document.addEventListener('DOMContentLoaded', () => {
    const deleteButton = document.querySelector('button[data-action="delete-post"]');
    if (!deleteButton) return;

    deleteButton.addEventListener('click', () => {
        const postID = deleteButton.dataset.postId;

        if (confirm("Are you sure you want to delete this post? This action cannot be undone.")) {
            fetch(`/view?id=${postID}`, {
                method: 'DELETE',
                headers: {
                    'csrf': document.getElementById('csrf').value,
                },
            })
                .then(response => {
                    if (response.ok) {
                        alert("Post deleted successfully.");
                        window.location.href = '/home';
                    } else {
                        alert("Failed to delete the post.");
                    }
                })
                .catch(error => {
                    console.error("Error:", error);
                    alert("An error occurred while deleting the post.");
                });
        }
    });
});

document.addEventListener('DOMContentLoaded', () => {
    // Initialize button styles based on data-vote-state
    document.querySelectorAll('[data-vote-state]').forEach(button => {
        const state = button.dataset.voteState;
        if (state === 'upvote' && button.id.startsWith('upvote-')) {
            button.classList.add('bg-blue-300');
        } else if (state === 'downvote' && button.id.startsWith('downvote-')) {
            button.classList.add('bg-red-200');
        }
    });

    // Function to send a vote and update the UI dynamically
    const sendVote = (id, isPost, voteType) => {
        const formData = new FormData();
        const csrf = document.getElementById('csrf');
        formData.append(isPost ? 'post_id' : 'comment_id', id);
        formData.append('vote_type', voteType);
        formData.append('csrf', csrf.value);
    
        fetch('/vote', {
            method: 'POST',
            body: formData,
        }).then(response => {
            if (response.ok) {
                response.json().then(data => {
                    // Update the vote counts dynamically
                    const upvoteButton = document.getElementById(`upvote-${isPost ? 'post' : 'comment'}-${id}`);
                    const downvoteButton = document.getElementById(`downvote-${isPost ? 'post' : 'comment'}-${id}`);
    
                    // Update the vote counts inside the <span> elements
                    upvoteButton.querySelector('.vote-count').textContent = data.upvotes;
                    downvoteButton.querySelector('.vote-count').textContent = data.downvotes;
    
                    // Reset button styles
                    upvoteButton.classList.remove('bg-blue-300');
                    downvoteButton.classList.remove('bg-red-200');
    
                    // Apply the new style based on the vote state
                    if (data.user_vote === 'upvote') {
                        upvoteButton.classList.add('bg-blue-300');
                    } else if (data.user_vote === 'downvote') {
                        downvoteButton.classList.add('bg-red-200');
                    }
                });
            } else {
                console.error('Failed to register vote');
                alert('Failed to register your vote. Please try again.');
            }
        }).catch(error => {
            console.error('Error:', error);
            alert('An error occurred while processing your vote.');
        });
    };

    // Add event listeners for upvote buttons
    document.querySelectorAll('[id^="upvote-"]').forEach(button => {
        button.addEventListener('click', () => {
            const id = button.id.split('-')[2];
            const isPost = button.id.includes('post');
            sendVote(id, isPost, 'upvote');
        });
    });

    // Add event listeners for downvote buttons
    document.querySelectorAll('[id^="downvote-"]').forEach(button => {
        button.addEventListener('click', () => {
            const id = button.id.split('-')[2];
            const isPost = button.id.includes('post');
            sendVote(id, isPost, 'downvote');
        });
    });
});